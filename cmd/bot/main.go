package main

import (
	"context"
	"log"
	"time"

	"github.com/HasanNugroho/coin-be/internal/bot"
	"github.com/HasanNugroho/coin-be/internal/bot/otp"
	"github.com/HasanNugroho/coin-be/internal/bot/session"
	"github.com/HasanNugroho/coin-be/internal/bot/vision"
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/core/database"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/dashboard"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	tele "gopkg.in/telebot.v4"
)

func main() {
	cfg := config.Load()

	// Database connection
	mongoClient, err := database.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(cfg.MongoDB)

	// Repositories
	userRepo := user.NewRepository(db)
	pocketRepo := pocket.NewRepository(db)
	transactionRepo := transaction.NewRepository(db)
	userPlatformRepo := user_platform.NewUserPlatformRepository(db)

	// Services
	transactionSvc := transaction.NewService(transactionRepo, pocketRepo, userPlatformRepo)
	dashboardSvc := dashboard.NewService(dashboard.NewRepository(db))

	// Bot components
	otpStore := otp.NewStore()
	sessionStore := session.NewStore()
	mailer := utils.NewMailer()
	visionParser := vision.NewReceiptParser(cfg.OpenAIKey, cfg.AIHost, cfg.AIModel)

	telegramSvc := bot.NewTelegramService(
		userRepo,
		transactionSvc,
		pocketRepo,
		dashboardSvc,
		otpStore,
		mailer,
		visionParser,
	)

	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	handler := bot.NewHandler(telegramSvc, sessionStore)
	handler.Register(b)

	log.Printf("Bot started as %s", b.Me.Username)
	b.Start()
}
