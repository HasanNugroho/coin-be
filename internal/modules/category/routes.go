package category

import (
	"github.com/HasanNugroho/coin-be/internal/core/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	protected := r.Group("")
	{
		protected.GET("", controller.ListCategories)
		protected.GET("/:id", controller.GetCategory)
		protected.GET("/type/:type", controller.ListCategoriesByType)
		protected.GET("/subcategories/:parent_id", controller.ListSubcategories)
	}

	admin := r.Group("")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("", controller.CreateCategory)
		admin.PUT("/:id", controller.UpdateCategory)
		admin.DELETE("/:id", controller.DeleteCategory)
	}
}
