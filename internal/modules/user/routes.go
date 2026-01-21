package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	protected := r.Group("")
	{
		protected.GET("/profile", controller.GetProfile)
		protected.PUT("/profile", controller.UpdateProfile)
		protected.GET("/financial-profile", controller.GetFinancialProfile)
		protected.POST("/financial-profile", controller.CreateFinancialProfile)
		protected.PUT("/financial-profile", controller.UpdateFinancialProfile)
		protected.DELETE("/financial-profile", controller.DeleteFinancialProfile)
	}

	admin := r.Group("")
	{
		admin.GET("", controller.ListUsers)
		admin.GET("/:id", controller.GetUser)
		admin.DELETE("/:id", controller.DeleteUser)
		admin.POST("/:id/disable", controller.DisableUser)
		admin.POST("/:id/enable", controller.EnableUser)
	}
}
