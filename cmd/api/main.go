package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/HasanNugroho/coin-be/internal/core/container"
	"github.com/HasanNugroho/coin-be/internal/core/middleware"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/auth"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/HasanNugroho/coin-be/internal/modules/health"
	_ "github.com/HasanNugroho/coin-be/docs"
)

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

	// User routes (protected)
	jwtManager := appContainer.Get("jwtManager").(*utils.JWTManager)
	userController := appContainer.Get("userController").(*user.Controller)
	userRoutes := api.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(jwtManager))
	user.RegisterRoutes(userRoutes, userController)

	// Health routes (public)
	healthController := appContainer.Get("healthController").(*health.Controller)
	health.RegisterRoutes(api.Group("/health"), healthController)

	log.Println("Server running on http://localhost:8080")
	log.Println("Swagger docs available at http://localhost:8080/swagger/index.html")
	r.Run(":8080")
}
