package auth

import (
	"github.com/HasanNugroho/coin-be/internal/core/middleware"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller, jwtManager *utils.JWTManager, db *mongo.Database) {
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.POST("/refresh", controller.RefreshToken)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtManager, db))
	protected.POST("/logout", controller.Logout)
	protected.GET("/me", controller.GetMe)
	protected.GET("/validate", controller.ValidateToken)
}
