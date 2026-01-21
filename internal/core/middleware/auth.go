package middleware

import (
	"net/http"
	"strings"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware(jwtManager *utils.JWTManager, db *mongo.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("missing authorization header"))
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("invalid authorization header"))
			ctx.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtManager.VerifyAccessToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("invalid token"))
			ctx.Abort()
			return
		}

		// Get user role from database
		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("invalid user id"))
			ctx.Abort()
			return
		}

		usersCollection := db.Collection("users")
		var user struct {
			Role string `bson:"role"`
		}
		err = usersCollection.FindOne(ctx, primitive.M{"_id": userID}).Decode(&user)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("user not found"))
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Set("role", user.Role)
		ctx.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(http.StatusForbidden, errors.NewErrorResponse("unable to verify role"))
			ctx.Abort()
			return
		}

		userRole := role.(string)
		if userRole != "admin" {
			ctx.JSON(http.StatusForbidden, errors.NewErrorResponse("admin access required"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
