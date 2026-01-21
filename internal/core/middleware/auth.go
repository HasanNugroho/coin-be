package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/pkg/errors"
)

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
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

		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Next()
	}
}

func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("unauthorized"))
			ctx.Abort()
			return
		}

		roles, exists := ctx.Get("roles")
		if !exists {
			ctx.JSON(http.StatusForbidden, errors.NewErrorResponse("unable to verify roles"))
			ctx.Abort()
			return
		}

		userRoles := roles.([]string)
		hasRequiredRole := false

		for _, userRole := range userRoles {
			for _, required := range requiredRoles {
				if userRole == required {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			ctx.JSON(http.StatusForbidden, errors.NewErrorResponse("insufficient permissions"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
