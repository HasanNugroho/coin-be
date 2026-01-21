package main

import (
	"log"

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
	"github.com/HasanNugroho/coin-be/internal/modules/allocation"
	"github.com/HasanNugroho/coin-be/internal/modules/auth"
	"github.com/HasanNugroho/coin-be/internal/modules/category"
	"github.com/HasanNugroho/coin-be/internal/modules/health"
	"github.com/HasanNugroho/coin-be/internal/modules/report"
	"github.com/HasanNugroho/coin-be/internal/modules/target"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
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
	health.Register(builder)
	category.Register(builder)
	allocation.Register(builder)
	transaction.Register(builder)
	target.Register(builder)
	report.Register(builder)

	appContainer := builder.Build()

	r := gin.Default()

	// Apply global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	// Auth routes (public)
	authController := appContainer.Get("authController").(*auth.Controller)
	authRoutes := api.Group("/auth")
	auth.RegisterRoutes(authRoutes, authController)

	// Get dependencies for middleware
	jwtManager := appContainer.Get("jwtManager").(*utils.JWTManager)
	mongoClient := appContainer.Get("mongo").(*mongo.Client)
	cfg := appContainer.Get("config").(*config.Config)
	db := mongoClient.Database(cfg.MongoDB)

	// User routes (protected)
	userController := appContainer.Get("userController").(*user.Controller)
	userRoutes := api.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	userRoutes.Use(middleware.AdminMiddleware())
	user.RegisterRoutes(userRoutes, userController)

	// Health routes (public)
	healthController := appContainer.Get("healthController").(*health.Controller)
	health.RegisterRoutes(api.Group("/health"), healthController)

	// Category routes (protected)
	categoryController := appContainer.Get("categoryController").(*category.Controller)
	categoryRoutes := api.Group("/categories")
	categoryRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	category.RegisterRoutes(categoryRoutes, categoryController)

	// Transaction routes (protected)
	transactionController := appContainer.Get("transactionController").(*transaction.Controller)
	transactionRoutes := api.Group("/transactions")
	transactionRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	transaction.RegisterRoutes(transactionRoutes, transactionController)

	// Allocation routes (protected)
	allocationController := appContainer.Get("allocationController").(*allocation.Controller)
	allocationRoutes := api.Group("/allocations")
	allocationRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	allocation.RegisterRoutes(allocationRoutes, allocationController)

	// Target routes (protected)
	targetController := appContainer.Get("targetController").(*target.Controller)
	targetRoutes := api.Group("/targets")
	targetRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	target.RegisterRoutes(targetRoutes, targetController)

	// Report routes (protected)
	reportController := appContainer.Get("reportController").(*report.Controller)
	reportRoutes := api.Group("/reports")
	reportRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
	report.RegisterRoutes(reportRoutes, reportController)

	log.Println("Server running on http://localhost:8080")
	log.Println("Swagger docs available at http://localhost:8080/swagger/index.html")
	r.Run(":8080")
}
