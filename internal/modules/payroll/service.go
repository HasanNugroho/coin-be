package payroll

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/allocation"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	payrollRepo      *Repository
	userRepo         *user.Repository
	userPlatformRepo *user_platform.UserPlatformRepository
	pocketRepo       *pocket.Repository
	allocationRepo   *allocation.Repository
	transactionRepo  *transaction.Repository
	balanceProcessor *transaction.BalanceProcessor
	db               *mongo.Database
}

func NewService(
	payrollRepo *Repository,
	userRepo *user.Repository,
	userPlatformRepo *user_platform.UserPlatformRepository,
	pocketRepo *pocket.Repository,
	allocationRepo *allocation.Repository,
	transactionRepo *transaction.Repository,
	balanceProcessor *transaction.BalanceProcessor,
	db *mongo.Database,
) *Service {
	return &Service{
		payrollRepo:      payrollRepo,
		userRepo:         userRepo,
		userPlatformRepo: userPlatformRepo,
		pocketRepo:       pocketRepo,
		allocationRepo:   allocationRepo,
		transactionRepo:  transactionRepo,
		balanceProcessor: balanceProcessor,
		db:               db,
	}
}

// ProcessDailyPayroll processes payroll for all eligible users on their salary day
func (s *Service) ProcessDailyPayroll(ctx context.Context) error {
	now := time.Now()
	today := now.Day()

	// Fetch all users with auto_input_payroll enabled
	users, err := s.userRepo.ListUsers(ctx, 1000, 0)
	if err != nil {
		log.Printf("failed to fetch users for payroll: %v", err)
		return err
	}

	successCount := 0
	failureCount := 0

	for _, u := range users {
		// Get user profile
		profile, err := s.userRepo.GetUserProfileByUserID(ctx, u.ID)
		if err != nil || profile == nil {
			continue
		}

		// Check if payroll is enabled and it's the salary day
		if !profile.AutoInputPayroll || profile.SalaryDay != today || profile.BaseSalary <= 0 {
			continue
		}

		// Check if payroll already processed today
		payrollRecord, err := s.payrollRepo.GetPayrollRecord(ctx, u.ID, now.Year(), int(now.Month()), today)
		if err != nil {
			log.Printf("error checking payroll record for user %s: %v", u.ID.Hex(), err)
			continue
		}

		if payrollRecord != nil {
			// Payroll already processed today
			continue
		}

		// Process payroll for this user
		err = s.processUserPayroll(ctx, u, profile)
		if err != nil {
			log.Printf("failed to process payroll for user %s: %v", u.ID.Hex(), err)
			failureCount++

			// Record failure
			errMsg := err.Error()
			s.payrollRepo.CreatePayrollRecord(ctx, &PayrollRecord{
				UserID: u.ID,
				Year:   now.Year(),
				Month:  int(now.Month()),
				Day:    today,
				Amount: profile.BaseSalary,
				Status: StatusFailed,
				Error:  &errMsg,
			})
		} else {
			log.Printf("successfully processed payroll for user %s", u.ID.Hex())
			successCount++

			// Record success
			s.payrollRepo.CreatePayrollRecord(ctx, &PayrollRecord{
				UserID: u.ID,
				Year:   now.Year(),
				Month:  int(now.Month()),
				Day:    today,
				Amount: profile.BaseSalary,
				Status: StatusSuccess,
			})
		}
	}

	log.Printf("payroll processing complete: %d success, %d failures", successCount, failureCount)
	return nil
}

// processUserPayroll processes payroll for a single user within a database transaction
func (s *Service) processUserPayroll(ctx context.Context, u *user.User, profile *user.UserProfile) error {
	// Validate default user platform
	if profile.DefaultUserPlatformID == nil {
		return errors.New("no default user platform configured")
	}

	defaultUserPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, *profile.DefaultUserPlatformID)
	if err != nil {
		return errors.New("default user platform not found")
	}

	if !defaultUserPlatform.IsActive {
		return errors.New("default user platform is not active")
	}

	if defaultUserPlatform.UserID != u.ID {
		return errors.New("default user platform does not belong to user")
	}

	// Get main pocket for income transaction
	mainPocket, err := s.getMainPocket(ctx, u.ID)
	if err != nil {
		return err
	}

	// Start database transaction
	session, err := s.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionCtx mongo.SessionContext) error {
		// Step 1: Create INCOME transaction
		incomeTransaction := &transaction.Transaction{
			UserID:         u.ID,
			Type:           string(transaction.TypeIncome),
			Amount:         profile.BaseSalary,
			PocketTo:       &mainPocket.ID,
			UserPlatformTo: &defaultUserPlatform.ID,
			Date:           time.Now(),
			Note:           stringPtr("Payroll auto-input"),
			Ref:            stringPtr("payroll_" + time.Now().Format("2006_01_02")),
		}

		// Persist income transaction
		if err := s.transactionRepo.CreateTransaction(sessionCtx, incomeTransaction); err != nil {
			return errors.New("failed to create income transaction: " + err.Error())
		}

		// Step 2: Update balances for income transaction
		if err := s.updateBalancesForIncome(sessionCtx, mainPocket, defaultUserPlatform, profile.BaseSalary); err != nil {
			return errors.New("failed to update balances for income: " + err.Error())
		}

		// Step 3: Execute allocations
		allocations, err := s.allocationRepo.GetActiveAllocationsByUserID(sessionCtx, u.ID)
		if err != nil {
			return errors.New("failed to fetch allocations: " + err.Error())
		}

		// Build all allocation transactions in memory
		allocationTransactions := make([]*transaction.Transaction, 0)
		remainingBalance := profile.BaseSalary

		for _, alloc := range allocations {
			if !alloc.IsActive {
				continue
			}

			// Calculate allocation amount
			var allocAmount float64
			if alloc.AllocationType == string(allocation.TypePercentage) {
				allocAmount = remainingBalance * (alloc.Nominal / 100)
			} else {
				allocAmount = alloc.Nominal
			}

			// Validate sufficient balance
			if allocAmount > remainingBalance {
				continue
			}

			// Validate target
			var targetPocket *primitive.ObjectID
			var targetUserPlatform *primitive.ObjectID

			if alloc.PocketID != nil {
				pocket, err := s.pocketRepo.GetPocketByID(sessionCtx, *alloc.PocketID)
				if err != nil || !pocket.IsActive || pocket.UserID != u.ID {
					continue
				}
				targetPocket = alloc.PocketID
			}

			if alloc.UserPlatformID != nil {
				userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(sessionCtx, *alloc.UserPlatformID)
				if err != nil || !userPlatform.IsActive || userPlatform.UserID != u.ID {
					continue
				}
				targetUserPlatform = alloc.UserPlatformID
			}

			if targetPocket == nil && targetUserPlatform == nil {
				continue
			}

			// Create allocation transaction
			allocTx := &transaction.Transaction{
				UserID:           u.ID,
				Type:             string(transaction.TypeTransfer),
				Amount:           allocAmount,
				PocketFrom:       &mainPocket.ID,
				PocketTo:         targetPocket,
				UserPlatformFrom: &defaultUserPlatform.ID,
				UserPlatformTo:   targetUserPlatform,
				Date:             time.Now(),
				Note:             stringPtr("Auto allocation - priority " + string(rune(alloc.Priority+'0'))),
				Ref:              stringPtr("alloc_" + alloc.ID.Hex()),
			}

			allocationTransactions = append(allocationTransactions, allocTx)
			remainingBalance -= allocAmount
		}

		// Step 4: Bulk insert all allocation transactions
		if len(allocationTransactions) > 0 {
			txInterfaces := make([]interface{}, len(allocationTransactions))
			for i, tx := range allocationTransactions {
				tx.ID = primitive.NewObjectID()
				tx.CreatedAt = time.Now()
				txInterfaces[i] = tx
			}

			col := s.db.Collection("transactions")
			if _, err := col.InsertMany(sessionCtx, txInterfaces); err != nil {
				return errors.New("failed to bulk insert allocation transactions: " + err.Error())
			}

			// Step 5: Update balances for all allocations
			for _, allocTx := range allocationTransactions {
				if err := s.updateBalancesForTransfer(sessionCtx, allocTx); err != nil {
					return errors.New("failed to update allocation balances: " + err.Error())
				}
			}
		}

		return nil
	})

	return err
}

// updateBalancesForIncome updates balances for income transaction
func (s *Service) updateBalancesForIncome(ctx context.Context, mainPocket *pocket.Pocket, userPlatform *user_platform.UserPlatform, amount float64) error {
	// Update main pocket balance
	mainPocket.Balance = utils.AddDecimal128(mainPocket.Balance, amount)
	if err := s.pocketRepo.UpdatePocket(ctx, mainPocket.ID, mainPocket); err != nil {
		return err
	}

	// Update user platform balance
	userPlatform.Balance = utils.AddDecimal128(userPlatform.Balance, amount)
	if err := s.userPlatformRepo.UpdateUserPlatform(ctx, userPlatform.ID, userPlatform); err != nil {
		return err
	}

	return nil
}

// updateBalancesForTransfer updates balances for transfer transaction
func (s *Service) updateBalancesForTransfer(ctx context.Context, tx *transaction.Transaction) error {
	// Decrease source balances
	if tx.PocketFrom != nil {
		pocket, err := s.pocketRepo.GetPocketByID(ctx, *tx.PocketFrom)
		if err != nil {
			return err
		}
		pocket.Balance = utils.AddDecimal128(pocket.Balance, -tx.Amount)
		if err := s.pocketRepo.UpdatePocket(ctx, pocket.ID, pocket); err != nil {
			return err
		}
	}

	if tx.UserPlatformFrom != nil {
		userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, *tx.UserPlatformFrom)
		if err != nil {
			return err
		}
		userPlatform.Balance = utils.AddDecimal128(userPlatform.Balance, -tx.Amount)
		if err := s.userPlatformRepo.UpdateUserPlatform(ctx, userPlatform.ID, userPlatform); err != nil {
			return err
		}
	}

	// Increase target balances
	if tx.PocketTo != nil {
		pocket, err := s.pocketRepo.GetPocketByID(ctx, *tx.PocketTo)
		if err != nil {
			return err
		}
		pocket.Balance = utils.AddDecimal128(pocket.Balance, tx.Amount)
		if err := s.pocketRepo.UpdatePocket(ctx, pocket.ID, pocket); err != nil {
			return err
		}
	}

	if tx.UserPlatformTo != nil {
		userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, *tx.UserPlatformTo)
		if err != nil {
			return err
		}
		userPlatform.Balance = utils.AddDecimal128(userPlatform.Balance, tx.Amount)
		if err := s.userPlatformRepo.UpdateUserPlatform(ctx, userPlatform.ID, userPlatform); err != nil {
			return err
		}
	}

	return nil
}

// getMainPocket retrieves the main pocket for a user
func (s *Service) getMainPocket(ctx context.Context, userID primitive.ObjectID) (*pocket.Pocket, error) {
	pockets, err := s.pocketRepo.GetPocketsByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, err
	}

	for _, p := range pockets {
		if p.Type == string(pocket.TypeMain) && p.IsActive {
			return p, nil
		}
	}

	return nil, errors.New("no active main pocket found")
}

// stringPtr converts a string to a pointer, returning nil for empty strings
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
