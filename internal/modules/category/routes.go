package category

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	r.POST("", controller.CreateCategory)
	r.GET("", controller.GetCategories)
	r.GET("/:id", controller.GetCategoryByID)
	r.PUT("/:id", controller.UpdateCategory)
	r.DELETE("/:id", controller.DeleteCategory)
}
