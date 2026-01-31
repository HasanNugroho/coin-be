package platform

import (
	"github.com/HasanNugroho/coin-be/internal/core/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	protected := r.Group("")
	{
		protected.GET("", controller.ListPlatforms)
		protected.GET("/active", controller.ListActivePlatforms)
		protected.GET("/:id", controller.GetPlatform)
		protected.GET("/type/:type", controller.ListPlatformsByType)
	}

	admin := r.Group("admin")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("", controller.CreatePlatform)
		admin.PUT("/:id", controller.UpdatePlatform)
		admin.DELETE("/:id", controller.DeletePlatform)
	}
}
