package user_platform

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	// User routes
	protected := r.Group("")
	{
		protected.POST("", controller.CreateUserPlatform)
		protected.GET("", controller.ListUserPlatforms)
		protected.GET("/dropdown", controller.ListUserPlatformsDropdown)
		protected.GET("/:id", controller.GetUserPlatform)
		protected.PUT("/:id", controller.UpdateUserPlatform)
		protected.DELETE("/:id", controller.DeleteUserPlatform)
	}
}
