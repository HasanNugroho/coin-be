package bot

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/HasanNugroho/coin-be/internal/bot/otp"
	"github.com/HasanNugroho/coin-be/internal/bot/vision"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/dashboard"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction/dto"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TelegramService struct {
	userRepo       *user.Repository
	transactionSvc *transaction.Service
	pocketRepo     *pocket.Repository
	platformRepo   *user_platform.UserPlatformRepository
	dashboardSvc   *dashboard.Service
	otpStore       *otp.Store
	mailer         utils.Mailer
	visionParser   *vision.ReceiptParser
}

func NewTelegramService(
	userRepo *user.Repository,
	transactionSvc *transaction.Service,
	pocketRepo *pocket.Repository,
	platformRepo *user_platform.UserPlatformRepository,
	dashboardSvc *dashboard.Service,
	otpStore *otp.Store,
	mailer utils.Mailer,
	visionParser *vision.ReceiptParser,
) *TelegramService {
	return &TelegramService{
		userRepo:       userRepo,
		transactionSvc: transactionSvc,
		pocketRepo:     pocketRepo,
		platformRepo:   platformRepo,
		dashboardSvc:   dashboardSvc,
		otpStore:       otpStore,
		mailer:         mailer,
		visionParser:   visionParser,
	}
}

func (s *TelegramService) FindUserByTelegramID(ctx context.Context, telegramID string) (*user.User, error) {
	return s.userRepo.FindByTelegramID(ctx, telegramID)
}

func (s *TelegramService) FindUserByEmail(ctx context.Context, email string) (*user.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *TelegramService) SendOTP(ctx context.Context, email string, telegramID int64) error {
	code, _ := generateOTP(6)
	s.otpStore.Set(email, otp.Entry{
		OTP:        code,
		TelegramID: telegramID,
		ExpiresAt:  time.Now().Add(5 * time.Minute),
	})

	subject := "Coin Bot - Your Verification Code"
	body := fmt.Sprintf("Your OTP for Telegram registration is: %s\nThis code will expire in 5 minutes.", code)

	return s.mailer.Send(ctx, email, subject, body)
}

func (s *TelegramService) VerifyOTP(ctx context.Context, email, code string, telegramID int64) (*user.User, error) {
	entry, ok := s.otpStore.Get(email)
	if !ok || entry.OTP != code || entry.TelegramID != telegramID {
		return nil, fmt.Errorf("invalid or expired OTP")
	}

	if err := s.userRepo.SetTelegramID(ctx, email, fmt.Sprintf("%d", telegramID)); err != nil {
		return nil, err
	}

	s.otpStore.Delete(email)
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *TelegramService) GetSummary(ctx context.Context, userID primitive.ObjectID, timeRange string) (*dashboard.DashboardSummary, error) {
	return s.dashboardSvc.GetDashboardSummary(ctx, userID.Hex(), dashboard.TimeRange(timeRange))
}

func (s *TelegramService) GetPockets(ctx context.Context, userID primitive.ObjectID) ([]*pocket.Pocket, error) {
	return s.pocketRepo.GetPocketsByUserIDDropdown(ctx, userID)
}

func (s *TelegramService) GetPlatforms(ctx context.Context, userID primitive.ObjectID) ([]*user_platform.UserPlatform, error) {
	return s.platformRepo.GetUserPlatformsByUserIDDropdown(ctx, userID)
}

func (s *TelegramService) CreateTransaction(ctx context.Context, userID primitive.ObjectID, txType string, amount float64, pocketID, platformID, note, date string) error {
	req := &dto.CreateTransactionRequest{
		Type:   txType,
		Amount: amount,
		Note:   note,
		Date:   date,
	}

	if txType == "income" {
		req.PocketToID = pocketID
		req.UserPlatformToID = platformID
	} else {
		req.PocketFromID = pocketID
		req.UserPlatformFromID = platformID
	}

	_, err := s.transactionSvc.CreateTransaction(ctx, userID.Hex(), req)
	return err
}

func (s *TelegramService) ParseReceiptImage(ctx context.Context, imageData []byte) (*vision.ParsedReceipt, error) {
	return s.visionParser.Parse(ctx, imageData)
}

func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[num.Int64()]
	}
	return string(result), nil
}
