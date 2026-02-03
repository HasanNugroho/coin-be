package payroll

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	payrollRepo      *Repository
	userRepo         *user.Repository
	userPlatformRepo *user_platform.UserPlatformRepository
	pocketRepo       *pocket.Repository
	transactionRepo  *transaction.Repository
	balanceProcessor *transaction.BalanceProcessor
	db               *mongo.Database
}

func NewService(
	payrollRepo *Repository,
	userRepo *user.Repository,
	userPlatformRepo *user_platform.UserPlatformRepository,
	pocketRepo *pocket.Repository,
	transactionRepo *transaction.Repository,
	balanceProcessor *transaction.BalanceProcessor,
	db *mongo.Database,
) *Service {
	return &Service{
		payrollRepo:      payrollRepo,
		userRepo:         userRepo,
		userPlatformRepo: userPlatformRepo,
		pocketRepo:       pocketRepo,
		transactionRepo:  transactionRepo,
		balanceProcessor: balanceProcessor,
		db:               db,
	}
}

// ProcessDailyPayroll processes payroll for all eligible users on their salary day
func (s *Service) ProcessDailyPayroll(ctx context.Context) error {
	now := time.Now()
	today := now.Day()

	// Fetch eligible users using aggregation pipeline
	eligibleUsers, err := s.getEligibleUsersForPayroll(ctx, today)
	if err != nil {
		log.Printf("failed to fetch eligible users for payroll: %v", err)
		return err
	}

	if len(eligibleUsers) == 0 {
		log.Printf("no eligible users for payroll processing today")
		return nil
	}

	// Extract user IDs for bulk payroll record check
	userIDs := make([]primitive.ObjectID, len(eligibleUsers))
	for i, userData := range eligibleUsers {
		userIDs[i] = userData.Profile.UserID
	}

	// Check existing payroll records in bulk
	existingRecords, err := s.payrollRepo.GetPayrollRecordByUserIDs(ctx, userIDs, now.Year(), int(now.Month()), today)
	if err != nil {
		log.Printf("error checking payroll records: %v", err)
		// Continue processing, don't fail entirely
		existingRecords = []*PayrollRecord{}
	}

	// Create a map for faster lookup
	existingRecordsMap := make(map[primitive.ObjectID]bool)
	for _, record := range existingRecords {
		existingRecordsMap[record.UserID] = true
	}

	successCount := 0
	failureCount := 0
	newPayrollRecords := make([]*PayrollRecord, 0, len(eligibleUsers))

	for _, userData := range eligibleUsers {
		// Check if payroll already processed today
		if existingRecordsMap[userData.Profile.UserID] {
			log.Printf("payroll already processed for user %s", userData.Profile.UserID.Hex())
			continue
		}

		// Process payroll for this user
		err = s.processUserPayroll(ctx, &userData.User, &userData.Profile)

		// Prepare payroll record
		record := &PayrollRecord{
			UserID: userData.Profile.UserID,
			Year:   now.Year(),
			Month:  int(now.Month()),
			Day:    today,
			Amount: userData.Profile.BaseSalary,
		}

		if err != nil {
			log.Printf("failed to process payroll for user %s: %v", userData.Profile.UserID.Hex(), err)
			failureCount++
			errMsg := err.Error()
			record.Status = StatusFailed
			record.Error = &errMsg
		} else {
			log.Printf("successfully processed payroll for user %s", userData.Profile.UserID.Hex())
			successCount++
			record.Status = StatusSuccess
		}

		newPayrollRecords = append(newPayrollRecords, record)
	}

	// Bulk insert all new payroll records
	if len(newPayrollRecords) > 0 {
		if err := s.payrollRepo.CreatePayrollRecordBulk(ctx, newPayrollRecords); err != nil {
			log.Printf("failed to bulk insert payroll records: %v", err)
			// Don't fail the entire process if record creation fails
		}
	}

	log.Printf("payroll processing complete: %d success, %d failures", successCount, failureCount)
	return nil
}

// UserPayrollData represents the joined data from user_profiles and users
type UserPayrollData struct {
	Profile user.UserProfile `bson:",inline"`
	User    user.User        `bson:"user"`
}

// getEligibleUsersForPayroll fetches users eligible for payroll using aggregation
func (s *Service) getEligibleUsersForPayroll(ctx context.Context, salaryDay int) ([]UserPayrollData, error) {
	pipeline := mongo.Pipeline{
		// Match profiles with the given salary day and auto_input_payroll enabled
		{{Key: "$match", Value: bson.D{
			{Key: "salary_day", Value: salaryDay},
			{Key: "auto_input_payroll", Value: true},
			{Key: "base_salary", Value: bson.D{{Key: "$gt", Value: 0}}},
			{Key: "is_active", Value: true},
		}}},
		// Lookup user data
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		// Unwind user array
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},
		// Filter only active users
		{{Key: "$match", Value: bson.D{
			{Key: "user.is_active", Value: true},
		}}},
	}

	profileCol := s.db.Collection("user_profiles")
	cursor, err := profileCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]UserPayrollData, 0)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
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
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Step 1: Create INCOME transaction
		incomeTransaction := &transaction.Transaction{
			UserID:           u.ID,
			Type:             string(transaction.TypeIncome),
			Amount:           profile.BaseSalary,
			PocketToID:       &mainPocket.ID,
			UserPlatformToID: &defaultUserPlatform.ID,
			Date:             time.Now(),
			Note:             stringPtr("Payroll auto-input"),
			Ref:              stringPtr("payroll_" + time.Now().Format("2006_01_02")),
		}

		// Persist income transaction
		if err := s.transactionRepo.CreateTransaction(sessionCtx, incomeTransaction); err != nil {
			session.AbortTransaction(sessionCtx)
			return errors.New("failed to create income transaction: " + err.Error())
		}

		// Step 2: Update balances for income transaction
		if err := s.updateBalancesForIncome(sessionCtx, mainPocket, defaultUserPlatform, profile.BaseSalary); err != nil {
			session.AbortTransaction(sessionCtx)
			return errors.New("failed to update balances for income: " + err.Error())
		}

		return session.CommitTransaction(sessionCtx)
	})

	return err
}

// balanceUpdate tracks a balance change for batch processing
type balanceUpdate struct {
	entityType string // "pocket" or "userPlatform"
	id         primitive.ObjectID
	delta      float64
}

// applyBalanceUpdates applies all balance changes in batched operations
func (s *Service) applyBalanceUpdates(ctx context.Context, updates []balanceUpdate) error {
	// Consolidate updates by entity
	pocketUpdates := make(map[primitive.ObjectID]float64)
	userPlatformUpdates := make(map[primitive.ObjectID]float64)

	for _, update := range updates {
		if update.entityType == "pocket" {
			pocketUpdates[update.id] += update.delta
		} else if update.entityType == "userPlatform" {
			userPlatformUpdates[update.id] += update.delta
		}
	}

	// Apply pocket updates
	for pocketID, delta := range pocketUpdates {
		pocket, err := s.pocketRepo.GetPocketByID(ctx, pocketID)
		if err != nil {
			return err
		}
		pocket.Balance = utils.AddDecimal128(pocket.Balance, delta)
		if err := s.pocketRepo.UpdatePocket(ctx, pocket.ID, pocket); err != nil {
			return err
		}
	}

	// Apply user platform updates
	for userPlatformID, delta := range userPlatformUpdates {
		userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, userPlatformID)
		if err != nil {
			return err
		}
		userPlatform.Balance = utils.AddDecimal128(userPlatform.Balance, delta)
		if err := s.userPlatformRepo.UpdateUserPlatform(ctx, userPlatform.ID, userPlatform); err != nil {
			return err
		}
	}

	return nil
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
