package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"

	_ "github.com/HasanNugroho/coin-be/docs"
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/core/container"
	"github.com/HasanNugroho/coin-be/internal/core/middleware"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/auth"
	"github.com/HasanNugroho/coin-be/internal/modules/category_template"
	"github.com/HasanNugroho/coin-be/internal/modules/platform"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket_template"
	"github.com/HasanNugroho/coin-be/internal/modules/reporting"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/HasanNugroho/coin-be/internal/modules/user_category"
)

// @title Coin Backend API
// @version 1.0
// @description A comprehensive financial management system with smart allocation engine, transaction tracking, and detailed reports
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	builder, _ := di.NewBuilder()

	// Core container
	ctn, err := container.BuildContainer()
	if err != nil {
		log.Fatal(err)
	}

	// Copy core definitions
	for name, def := range ctn.Definitions() {
		def.Name = name
		builder.Add(def)
	}

	// Register modules
	auth.Register(builder)
	user.Register(builder)
	category_template.Register(builder)
	user_category.Register(builder)
	platform.Register(builder)
	pocket_template.Register(builder)
	pocket.Register(builder)
	transaction.Register(builder)
	reporting.Register(builder)

	appContainer := builder.Build()

	r := gin.Default()

	// Apply global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	// Get dependencies for middleware
	jwtManager := appContainer.Get("jwtManager").(*utils.JWTManager)
	mongoClient := appContainer.Get("mongo").(*mongo.Client)
	cfg := appContainer.Get("config").(*config.Config)
	db := mongoClient.Database(cfg.MongoDB)

	// Auth routes (public)
	authController := appContainer.Get("authController").(*auth.Controller)
	authRoutes := api.Group("/v1/auth")
	auth.RegisterRoutes(authRoutes, authController, jwtManager, db)

	// User routes (protected)
	userController := appContainer.Get("userController").(*user.Controller)
	userRoutes := api.Group("/v1/users")
	userRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	user.RegisterRoutes(userRoutes, userController)

	// Category Template routes (protected)
	categoryTemplateController := appContainer.Get("categoryTemplateController").(*category_template.Controller)
	categoryTemplateRoutes := api.Group("/v1/category-templates")
	categoryTemplateRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	category_template.RegisterRoutes(categoryTemplateRoutes, categoryTemplateController)

	// User Category routes (protected)
	userCategoryController := appContainer.Get("userCategoryController").(*user_category.Controller)
	userCategoryRoutes := api.Group("/v1/user-categories")
	userCategoryRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	user_category.RegisterRoutes(userCategoryRoutes, userCategoryController)

	// Platform routes (protected)
	platformController := appContainer.Get("platformController").(*platform.Controller)
	platformRoutes := api.Group("/v1/platforms")
	platformRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	platform.RegisterRoutes(platformRoutes, platformController)

	// Pocket Template routes (protected, admin only)
	pocketTemplateController := appContainer.Get("pocketTemplateController").(*pocket_template.Controller)
	pocketTemplateRoutes := api.Group("/v1/pocket-templates")
	pocketTemplateRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	pocket_template.RegisterRoutes(pocketTemplateRoutes, pocketTemplateController)

	// Pocket routes (protected)
	pocketController := appContainer.Get("pocketController").(*pocket.Controller)
	pocketRoutes := api.Group("/v1/pockets")
	pocketRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	pocket.RegisterRoutes(pocketRoutes, pocketController)

	// Transaction routes (protected)
	transactionController := appContainer.Get("transactionController").(*transaction.Controller)
	transactionRoutes := api.Group("/v1/transactions")
	transactionRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	transaction.RegisterRoutes(transactionRoutes, transactionController)

	// Reporting routes (protected)
	dashboardController := appContainer.Get("dashboardController").(*reporting.DashboardController)
	reportingRoutes := api.Group("/v1/dashboard")
	reportingRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	reporting.RegisterRoutes(reportingRoutes, dashboardController)

	// Start cron scheduler for reporting jobs
	cronScheduler := appContainer.Get("reportingCronScheduler").(*reporting.CronScheduler)
	if err := cronScheduler.Start(); err != nil {
		log.Fatalf("Failed to start cron scheduler: %v", err)
	}
	log.Println("[Main] Reporting cron scheduler started")

	// Graceful shutdown handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("[Main] Shutdown signal received: %v", sig)
		cronScheduler.Stop()
		log.Println("[Main] Cron scheduler stopped")
		os.Exit(0)
	}()

	log.Println("Server running on http://localhost:8080")
	log.Println("Swagger docs available at http://localhost:8080/swagger/index.html")
	log.Println("Dashboard API available at http://localhost:8080/api/v1/dashboard")
	r.Run(":8080")
}
