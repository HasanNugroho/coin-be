package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		ctx.Next()

		duration := time.Since(startTime)
		statusCode := ctx.Writer.Status()

		log.Printf("[%s] %s %s - %d (%v)", method, path, ctx.Request.RemoteAddr, statusCode, duration)
	}
}
