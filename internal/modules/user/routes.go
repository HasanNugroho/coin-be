package user

import (
	"github.com/HasanNugroho/coin-be/internal/core/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	protected := r.Group("")
	{
		protected.GET("/profile", controller.GetProfile)
		protected.PUT("/profile", controller.UpdateProfile)
	}

	admin := r.Group("")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("", controller.ListUsers)
		admin.GET("/:id", controller.GetUser)
		admin.DELETE("/:id", controller.DeleteUser)
		admin.POST("/:id/disable", controller.DisableUser)
		admin.POST("/:id/enable", controller.EnableUser)
	}
}
