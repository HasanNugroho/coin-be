package user_category

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	protected := r.Group("")
	{
		protected.GET("", controller.GetUserCategories)
		protected.GET("/:id", controller.GetUserCategoryByID)
		protected.POST("", controller.CreateUserCategory)
		protected.PUT("/:id", controller.UpdateUserCategory)
		protected.DELETE("/:id", controller.DeleteUserCategory)
	}
}
