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

		admin.POST("/roles", controller.CreateRole)
		admin.GET("/roles", controller.ListRoles)
		admin.GET("/roles/:id", controller.GetRole)

		admin.POST("/:id/roles", controller.AssignRoleToUser)
		admin.GET("/:id/roles", controller.GetUserRoles)
		admin.DELETE("/:id/roles/:role_id", controller.RemoveRoleFromUser)
	}
}
