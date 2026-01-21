package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HasanNugroho/coin-be/pkg/errors"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				ctx.JSON(http.StatusInternalServerError, errors.NewErrorResponse("internal server error"))
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
