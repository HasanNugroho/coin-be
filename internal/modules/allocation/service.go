package allocation

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/allocation/dto"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo             *Repository
	pocketRepo       *pocket.Repository
	userPlatformRepo *user_platform.UserPlatformRepository
	userRepo         *user.Repository
	transactionRepo  *transaction.Repository
	db               *mongo.Database
}

func NewService(r *Repository, pr *pocket.Repository, upr *user_platform.UserPlatformRepository, ur *user.Repository, tr *transaction.Repository, db *mongo.Database) *Service {
	return &Service{
		repo:             r,
		pocketRepo:       pr,
		userPlatformRepo: upr,
		userRepo:         ur,
		transactionRepo:  tr,
		db:               db,
	}
}

func (s *Service) CreateAllocation(ctx context.Context, userID string, req *dto.CreateAllocationRequest) (*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	if !IsValidAllocationType(req.AllocationType) {
		return nil, errors.New("invalid allocation type")
	}

	if !IsValidPriority(req.Priority) {
		return nil, errors.New("invalid priority")
	}

	// Validate that at least one target is provided
	if req.PocketID == "" && req.UserPlatformID == "" {
		return nil, errors.New("either pocket_id or user_platform_id must be provided")
	}

	var pocketID *primitive.ObjectID
	var userPlatformID *primitive.ObjectID

	// Validate pocket if provided
	if req.PocketID != "" {
		pocketObjID, err := primitive.ObjectIDFromHex(req.PocketID)
		if err != nil {
			return nil, errors.New("invalid pocket id")
		}

		pocket, err := s.pocketRepo.GetPocketByID(ctx, pocketObjID)
		if err != nil {
			return nil, errors.New("pocket not found")
		}

		if pocket.UserID != userObjID {
			return nil, errors.New("unauthorized: pocket does not belong to user")
		}

		if !pocket.IsActive {
			return nil, errors.New("pocket is not active")
		}

		pocketID = &pocketObjID
	}

	// Validate user platform if provided
	if req.UserPlatformID != "" {
		userPlatformObjID, err := primitive.ObjectIDFromHex(req.UserPlatformID)
		if err != nil {
			return nil, errors.New("invalid user platform id")
		}

		userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, userPlatformObjID)
		if err != nil {
			return nil, errors.New("user platform not found")
		}

		if userPlatform.UserID != userObjID {
			return nil, errors.New("unauthorized: user platform does not belong to user")
		}

		if !userPlatform.IsActive {
			return nil, errors.New("user platform is not active")
		}

		userPlatformID = &userPlatformObjID
	}

	// Validate nominal based on allocation type
	if req.AllocationType == string(TypePercentage) && req.Nominal > 100 {
		return nil, errors.New("percentage cannot exceed 100")
	}

	allocation := &Allocation{
		UserID:         userObjID,
		PocketID:       pocketID,
		UserPlatformID: userPlatformID,
		Priority:       req.Priority,
		AllocationType: req.AllocationType,
		Nominal:        req.Nominal,
		IsActive:       true,
		ExecuteDay:     req.ExecuteDay,
	}

	err = s.repo.CreateAllocation(ctx, allocation)
	if err != nil {
		return nil, err
	}

	return allocation, nil
}

func (s *Service) GetAllocationByID(ctx context.Context, userID string, allocationID string) (*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	allocationObjID, err := primitive.ObjectIDFromHex(allocationID)
	if err != nil {
		return nil, errors.New("invalid allocation id")
	}

	allocation, err := s.repo.GetAllocationByID(ctx, allocationObjID)
	if err != nil {
		return nil, err
	}

	if allocation.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	return allocation, nil
}

func (s *Service) ListAllocations(ctx context.Context, userID string) ([]*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetAllocationsByUserID(ctx, userObjID)
}

func (s *Service) GetActiveAllocations(ctx context.Context, userID string) ([]*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetActiveAllocationsByUserID(ctx, userObjID)
}

func (s *Service) UpdateAllocation(ctx context.Context, userID string, allocationID string, req *dto.UpdateAllocationRequest) (*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	allocationObjID, err := primitive.ObjectIDFromHex(allocationID)
	if err != nil {
		return nil, errors.New("invalid allocation id")
	}

	allocation, err := s.repo.GetAllocationByID(ctx, allocationObjID)
	if err != nil {
		return nil, err
	}

	if allocation.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	// Update pocket if provided
	if req.PocketID != "" {
		pocketObjID, err := primitive.ObjectIDFromHex(req.PocketID)
		if err != nil {
			return nil, errors.New("invalid pocket id")
		}

		pocket, err := s.pocketRepo.GetPocketByID(ctx, pocketObjID)
		if err != nil {
			return nil, errors.New("pocket not found")
		}

		if pocket.UserID != userObjID {
			return nil, errors.New("unauthorized: pocket does not belong to user")
		}

		allocation.PocketID = &pocketObjID
	}

	// Update user platform if provided
	if req.UserPlatformID != "" {
		userPlatformObjID, err := primitive.ObjectIDFromHex(req.UserPlatformID)
		if err != nil {
			return nil, errors.New("invalid user platform id")
		}

		userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, userPlatformObjID)
		if err != nil {
			return nil, errors.New("user platform not found")
		}

		if userPlatform.UserID != userObjID {
			return nil, errors.New("unauthorized: user platform does not belong to user")
		}

		allocation.UserPlatformID = &userPlatformObjID
	}

	if req.Priority != nil {
		if !IsValidPriority(*req.Priority) {
			return nil, errors.New("invalid priority")
		}
		allocation.Priority = *req.Priority
	}

	if req.AllocationType != "" {
		if !IsValidAllocationType(req.AllocationType) {
			return nil, errors.New("invalid allocation type")
		}
		allocation.AllocationType = req.AllocationType
	}

	if req.Nominal != nil {
		if allocation.AllocationType == string(TypePercentage) && *req.Nominal > 100 {
			return nil, errors.New("percentage cannot exceed 100")
		}
		allocation.Nominal = *req.Nominal
	}

	if req.IsActive != nil {
		allocation.IsActive = *req.IsActive
	}

	if req.ExecuteDay != nil {
		allocation.ExecuteDay = req.ExecuteDay
	}

	err = s.repo.UpdateAllocation(ctx, allocationObjID, allocation)
	if err != nil {
		return nil, err
	}

	return allocation, nil
}

func (s *Service) DeleteAllocation(ctx context.Context, userID string, allocationID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	allocationObjID, err := primitive.ObjectIDFromHex(allocationID)
	if err != nil {
		return errors.New("invalid allocation id")
	}

	allocation, err := s.repo.GetAllocationByID(ctx, allocationObjID)
	if err != nil {
		return err
	}

	if allocation.UserID != userObjID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteAllocation(ctx, allocationObjID)
}

// ProcessDailyAllocations processes all allocations scheduled for execution on the current day
// Automatically handles month-end logic: allocations scheduled for days 29-31 will execute
// on the last day of shorter months (e.g., day 31 executes on Feb 28/29, Apr 30, etc.)
func (s *Service) ProcessDailyAllocations(ctx context.Context) error {
	jakartaLoc := getJakartaLocation()
	now := time.Now().In(jakartaLoc)
	currentDay := now.Day()
	lastDayOfMonth := getLastDayOfMonth(now)

	log.Printf("Processing allocations for day %d (last day of month: %d)", currentDay, lastDayOfMonth)

	// Fetch allocations with automatic overflow handling
	allocationsData, err := s.repo.GetAllocationsByExecuteDayWithOverflow(ctx, currentDay, lastDayOfMonth)
	if err != nil {
		log.Printf("failed to fetch allocations for execution on day %d: %v", currentDay, err)
		return err
	}

	if len(allocationsData) == 0 {
		log.Printf("no allocations to process for day %d", currentDay)
		return nil
	}

	log.Printf("found %d allocations to process", len(allocationsData))
	return s.processAllocations(ctx, allocationsData, currentDay)
}

// processAllocations handles the execution of a batch of allocations
func (s *Service) processAllocations(ctx context.Context, allocationsData []map[string]interface{}, executeDay int) error {
	if len(allocationsData) == 0 {
		log.Printf("no allocations scheduled for execution on day %d", executeDay)
		return nil
	}

	successCount := 0
	failureCount := 0

	for _, allocData := range allocationsData {
		allocationID := allocData["_id"].(primitive.ObjectID)
		userID := allocData["user_id"].(primitive.ObjectID)

		err := s.processAllocationExecution(ctx, userID, allocationID, allocData)
		if err != nil {
			log.Printf("failed to process allocation %s for user %s: %v", allocationID.Hex(), userID.Hex(), err)
			failureCount++
		} else {
			log.Printf("successfully processed allocation %s for user %s", allocationID.Hex(), userID.Hex())
			successCount++
		}
	}

	log.Printf("allocation processing complete for day %d: %d success, %d failures", executeDay, successCount, failureCount)
	return nil
}

// processAllocationExecution executes a single allocation within a database transaction
func (s *Service) processAllocationExecution(ctx context.Context, userID primitive.ObjectID, allocationID primitive.ObjectID, allocData map[string]interface{}) error {
	allocation, err := s.repo.GetAllocationByID(ctx, allocationID)
	if err != nil {
		return err
	}

	if !allocation.IsActive {
		return errors.New("allocation is not active")
	}

	userProfileInterface := allocData["user_profile"]

	userProfileMap := userProfileInterface.(map[string]interface{})

	defaultUserPlatformID := userProfileMap["default_user_platform_id"]
	if defaultUserPlatformID == nil {
		return errors.New("no default user platform configured")
	}

	defaultUserPlatformObjID, ok := defaultUserPlatformID.(primitive.ObjectID)
	if !ok {
		return errors.New("invalid default user platform id")
	}

	defaultUserPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, defaultUserPlatformObjID)
	if err != nil {
		return errors.New("default user platform not found")
	}

	if !defaultUserPlatform.IsActive {
		return errors.New("default user platform is not active")
	}

	mainPocket, err := s.getMainPocket(ctx, userID)
	if err != nil {
		return err
	}

	session, err := s.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionCtx mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		var targetPocketID *primitive.ObjectID
		var targetUserPlatformID *primitive.ObjectID

		if allocation.PocketID != nil {
			pocket, err := s.pocketRepo.GetPocketByID(sessionCtx, *allocation.PocketID)
			if err != nil || !pocket.IsActive || pocket.UserID != userID {
				session.AbortTransaction(sessionCtx)
				return errors.New("target pocket is invalid")
			}
			targetPocketID = allocation.PocketID
		}

		if allocation.UserPlatformID != nil {
			userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(sessionCtx, *allocation.UserPlatformID)
			if err != nil || !userPlatform.IsActive || userPlatform.UserID != userID {
				session.AbortTransaction(sessionCtx)
				return errors.New("target user platform is invalid")
			}
			targetUserPlatformID = allocation.UserPlatformID
		}

		if targetPocketID == nil && targetUserPlatformID == nil {
			session.AbortTransaction(sessionCtx)
			return errors.New("no valid target for allocation")
		}

		allocAmount := allocation.Nominal

		allocTx := &transaction.Transaction{
			UserID:             userID,
			Type:               string(transaction.TypeTransfer),
			Amount:             allocAmount,
			PocketFromID:       &mainPocket.ID,
			PocketToID:         targetPocketID,
			UserPlatformFromID: &defaultUserPlatform.ID,
			UserPlatformToID:   targetUserPlatformID,
			Date:               time.Now(),
			Note:               stringPtr("Scheduled allocation execution"),
			Ref:                stringPtr("alloc_exec_" + allocationID.Hex()),
		}

		if err := s.transactionRepo.CreateTransaction(sessionCtx, allocTx); err != nil {
			session.AbortTransaction(sessionCtx)
			return errors.New("failed to create allocation transaction: " + err.Error())
		}

		balanceUpdates := make([]balanceUpdate, 0, 4)
		balanceUpdates = append(balanceUpdates,
			balanceUpdate{entityType: "pocket", id: mainPocket.ID, delta: -allocAmount},
			balanceUpdate{entityType: "userPlatform", id: defaultUserPlatform.ID, delta: -allocAmount},
		)

		if targetPocketID != nil {
			balanceUpdates = append(balanceUpdates, balanceUpdate{entityType: "pocket", id: *targetPocketID, delta: allocAmount})
		}
		if targetUserPlatformID != nil {
			balanceUpdates = append(balanceUpdates, balanceUpdate{entityType: "userPlatform", id: *targetUserPlatformID, delta: allocAmount})
		}

		if err := s.applyBalanceUpdates(sessionCtx, balanceUpdates); err != nil {
			session.AbortTransaction(sessionCtx)
			return errors.New("failed to update balances: " + err.Error())
		}

		return session.CommitTransaction(sessionCtx)
	})

	return err
}

// balanceUpdate tracks a balance change for batch processing
type balanceUpdate struct {
	entityType string
	id         primitive.ObjectID
	delta      float64
}

// applyBalanceUpdates applies all balance changes in batched operations
func (s *Service) applyBalanceUpdates(ctx context.Context, updates []balanceUpdate) error {
	pocketUpdates := make(map[primitive.ObjectID]float64)
	userPlatformUpdates := make(map[primitive.ObjectID]float64)

	for _, update := range updates {
		if update.entityType == "pocket" {
			pocketUpdates[update.id] += update.delta
		} else if update.entityType == "userPlatform" {
			userPlatformUpdates[update.id] += update.delta
		}
	}

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

// getLastDayOfMonth returns the last day of the given month
func getLastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}
