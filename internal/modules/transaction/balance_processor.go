package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BalanceProcessor centralizes all balance update logic.
// This is the single source of truth for balance changes.
// No balance update is allowed outside this processor.
type BalanceProcessor struct {
	pocketRepo       *pocket.Repository
	userPlatformRepo *user_platform.UserPlatformRepository
}

func NewBalanceProcessor(pr *pocket.Repository, upr *user_platform.UserPlatformRepository) *BalanceProcessor {
	return &BalanceProcessor{
		pocketRepo:       pr,
		userPlatformRepo: upr,
	}
}

// ProcessTransaction applies balance changes based on transaction type.
// This function enforces strict balance rules:
// - Income: increases pocket_to and user_platform_to
// - Expense: decreases pocket_from and user_platform_from
// - Transfer (pocket-to-pocket): reallocates between pockets only
// - Transfer (platform-to-platform): moves between user platforms only
// - Transfer (platform+pocket): moves between platforms and reassigns pockets
func (bp *BalanceProcessor) ProcessTransaction(
	ctx context.Context,
	txType string,
	amount float64,
	pocketFrom, pocketTo *primitive.ObjectID,
	userPlatformFrom, userPlatformTo *primitive.ObjectID,
) error {
	switch txType {
	case string(TypeIncome):
		return bp.processIncome(ctx, amount, pocketTo, userPlatformTo)

	case string(TypeExpense):
		return bp.processExpense(ctx, amount, pocketFrom, userPlatformFrom)

	case string(TypeTransfer):
		return bp.processTransfer(ctx, amount, pocketFrom, pocketTo, userPlatformFrom, userPlatformTo)

	default:
		return errors.New("invalid transaction type")
	}
}

// RevertTransaction reverses balance changes.
func (bp *BalanceProcessor) RevertTransaction(
	ctx context.Context,
	txType string,
	amount float64,
	pocketFrom, pocketTo *primitive.ObjectID,
	userPlatformFrom, userPlatformTo *primitive.ObjectID,
) error {
	switch txType {
	case string(TypeIncome):
		// Reverse income: decrease balances
		return bp.processExpense(ctx, amount, pocketTo, userPlatformTo)

	case string(TypeExpense):
		// Reverse expense: increase balances
		return bp.processIncome(ctx, amount, pocketFrom, userPlatformFrom)

	case string(TypeTransfer):
		// Reverse transfer: swap source and destination
		return bp.processTransfer(ctx, amount, pocketTo, pocketFrom, userPlatformTo, userPlatformFrom)

	default:
		return errors.New("invalid transaction type")
	}
}

// processIncome increases pocket_to and user_platform_to balances.
// Income brings money into the system.
func (bp *BalanceProcessor) processIncome(
	ctx context.Context,
	amount float64,
	pocketTo *primitive.ObjectID,
	userPlatformTo *primitive.ObjectID,
) error {
	// Update pocket balance if provided
	if pocketTo != nil {
		pocket, err := bp.pocketRepo.GetPocketByID(ctx, *pocketTo)
		if err != nil {
			return err
		}
		pocket.Balance = utils.AddDecimal128(pocket.Balance, amount)
		pocket.LastUseAt = time.Now()
		if err := bp.pocketRepo.UpdatePocket(ctx, *pocketTo, pocket); err != nil {
			return err
		}
	}

	// Update user platform balance if provided
	if userPlatformTo != nil {
		userPlatform, err := bp.userPlatformRepo.GetUserPlatformByID(ctx, *userPlatformTo)
		if err != nil {
			return err
		}
		userPlatform.Balance = utils.AddDecimal128(userPlatform.Balance, amount)
		userPlatform.LastUseAt = time.Now()
		if err := bp.userPlatformRepo.UpdateUserPlatform(ctx, *userPlatformTo, userPlatform); err != nil {
			return err
		}
	}

	return nil
}

// processExpense decreases pocket_from and user_platform_from balances.
// Expense removes money from the system.
func (bp *BalanceProcessor) processExpense(
	ctx context.Context,
	amount float64,
	pocketFrom *primitive.ObjectID,
	userPlatformFrom *primitive.ObjectID,
) error {
	// Update pocket balance if provided
	if pocketFrom != nil {
		pocket, err := bp.pocketRepo.GetPocketByID(ctx, *pocketFrom)
		if err != nil {
			return err
		}
		pocket.Balance = utils.AddDecimal128(pocket.Balance, -amount)
		pocket.LastUseAt = time.Now()
		if err := bp.pocketRepo.UpdatePocket(ctx, *pocketFrom, pocket); err != nil {
			return err
		}
	}

	// Update user platform balance if provided
	if userPlatformFrom != nil {
		userPlatform, err := bp.userPlatformRepo.GetUserPlatformByID(ctx, *userPlatformFrom)
		if err != nil {
			return err
		}
		userPlatform.Balance = utils.AddDecimal128(userPlatform.Balance, -amount)
		userPlatform.LastUseAt = time.Now()
		if err := bp.userPlatformRepo.UpdateUserPlatform(ctx, *userPlatformFrom, userPlatform); err != nil {
			return err
		}
	}

	return nil
}

// processTransfer handles three scenarios:
// 1. Pocket-to-pocket: reallocates between pockets only (no platform balance change)
// 2. Platform-to-platform: moves between user platforms only (no pocket balance change)
// 3. Platform+pocket: moves between platforms and reassigns pockets
func (bp *BalanceProcessor) processTransfer(
	ctx context.Context,
	amount float64,
	pocketFrom, pocketTo *primitive.ObjectID,
	userPlatformFrom, userPlatformTo *primitive.ObjectID,
) error {
	// Scenario 1: Pocket-to-pocket transfer (no platform balance change)
	if pocketFrom != nil && pocketTo != nil && userPlatformFrom == nil && userPlatformTo == nil {
		return bp.transferBetweenPockets(ctx, amount, pocketFrom, pocketTo)
	}

	// Scenario 2: Platform-to-platform transfer (no pocket balance change)
	if userPlatformFrom != nil && userPlatformTo != nil && pocketFrom == nil && pocketTo == nil {
		return bp.transferBetweenUserPlatforms(ctx, amount, userPlatformFrom, userPlatformTo)
	}

	// Scenario 3: Platform+pocket transfer (both platforms and pockets involved)
	if pocketFrom != nil && pocketTo != nil && userPlatformFrom != nil && userPlatformTo != nil {
		return bp.transferBetweenPlatformsWithPockets(ctx, amount, pocketFrom, pocketTo, userPlatformFrom, userPlatformTo)
	}

	return errors.New("invalid transfer combination: must specify either (pocket_from + pocket_to) or (user_platform_from + user_platform_to) or both pairs")
}

// transferBetweenPockets reallocates money between two pockets.
// Platform balance is unchanged.
func (bp *BalanceProcessor) transferBetweenPockets(
	ctx context.Context,
	amount float64,
	pocketFrom, pocketTo *primitive.ObjectID,
) error {
	// Decrease source pocket
	pocketFromData, err := bp.pocketRepo.GetPocketByID(ctx, *pocketFrom)
	if err != nil {
		return err
	}
	pocketFromData.Balance = utils.AddDecimal128(pocketFromData.Balance, -amount)
	pocketFromData.LastUseAt = time.Now()
	if err := bp.pocketRepo.UpdatePocket(ctx, *pocketFrom, pocketFromData); err != nil {
		return err
	}

	// Increase destination pocket
	pocketToData, err := bp.pocketRepo.GetPocketByID(ctx, *pocketTo)
	if err != nil {
		return err
	}
	pocketToData.Balance = utils.AddDecimal128(pocketToData.Balance, amount)
	pocketToData.LastUseAt = time.Now()
	if err := bp.pocketRepo.UpdatePocket(ctx, *pocketTo, pocketToData); err != nil {
		return err
	}

	return nil
}

// transferBetweenUserPlatforms moves money between user platforms.
// Pocket balance is unchanged.
func (bp *BalanceProcessor) transferBetweenUserPlatforms(
	ctx context.Context,
	amount float64,
	userPlatformFrom, userPlatformTo *primitive.ObjectID,
) error {
	// Decrease source user platform
	userPlatformFromData, err := bp.userPlatformRepo.GetUserPlatformByID(ctx, *userPlatformFrom)
	if err != nil {
		return err
	}
	userPlatformFromData.Balance = utils.AddDecimal128(userPlatformFromData.Balance, -amount)
	userPlatformFromData.LastUseAt = time.Now()
	if err := bp.userPlatformRepo.UpdateUserPlatform(ctx, *userPlatformFrom, userPlatformFromData); err != nil {
		return err
	}

	// Increase destination user platform
	userPlatformToData, err := bp.userPlatformRepo.GetUserPlatformByID(ctx, *userPlatformTo)
	if err != nil {
		return err
	}
	userPlatformToData.Balance = utils.AddDecimal128(userPlatformToData.Balance, amount)
	userPlatformToData.LastUseAt = time.Now()
	if err := bp.userPlatformRepo.UpdateUserPlatform(ctx, *userPlatformTo, userPlatformToData); err != nil {
		return err
	}

	return nil
}

// transferBetweenPlatformsWithPockets moves money between platforms and reassigns pockets.
// Both platform and pocket balances change.
func (bp *BalanceProcessor) transferBetweenPlatformsWithPockets(
	ctx context.Context,
	amount float64,
	pocketFrom, pocketTo *primitive.ObjectID,
	userPlatformFrom, userPlatformTo *primitive.ObjectID,
) error {
	// Decrease source pocket
	pocketFromData, err := bp.pocketRepo.GetPocketByID(ctx, *pocketFrom)
	if err != nil {
		return err
	}
	pocketFromData.Balance = utils.AddDecimal128(pocketFromData.Balance, -amount)
	pocketFromData.LastUseAt = time.Now()
	if err := bp.pocketRepo.UpdatePocket(ctx, *pocketFrom, pocketFromData); err != nil {
		return err
	}

	// Increase destination pocket
	pocketToData, err := bp.pocketRepo.GetPocketByID(ctx, *pocketTo)
	if err != nil {
		return err
	}
	pocketToData.Balance = utils.AddDecimal128(pocketToData.Balance, amount)
	pocketToData.LastUseAt = time.Now()
	if err := bp.pocketRepo.UpdatePocket(ctx, *pocketTo, pocketToData); err != nil {
		return err
	}

	// Decrease source user platform
	userPlatformFromData, err := bp.userPlatformRepo.GetUserPlatformByID(ctx, *userPlatformFrom)
	if err != nil {
		return err
	}
	userPlatformFromData.Balance = utils.AddDecimal128(userPlatformFromData.Balance, -amount)
	userPlatformFromData.LastUseAt = time.Now()
	if err := bp.userPlatformRepo.UpdateUserPlatform(ctx, *userPlatformFrom, userPlatformFromData); err != nil {
		return err
	}

	// Increase destination user platform
	userPlatformToData, err := bp.userPlatformRepo.GetUserPlatformByID(ctx, *userPlatformTo)
	if err != nil {
		return err
	}
	userPlatformToData.Balance = utils.AddDecimal128(userPlatformToData.Balance, amount)
	userPlatformToData.LastUseAt = time.Now()
	if err := bp.userPlatformRepo.UpdateUserPlatform(ctx, *userPlatformTo, userPlatformToData); err != nil {
		return err
	}

	return nil
}
