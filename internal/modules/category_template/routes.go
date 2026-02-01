package category_template

import (
	"github.com/HasanNugroho/coin-be/internal/core/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	admin := r.Group("")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("", controller.FindAll)
		admin.GET("/parent", controller.FindAllParent)
		admin.GET("/:id", controller.GetCategoryTemplateByID)
		admin.POST("", controller.CreateCategoryTemplate)
		admin.PUT("/:id", controller.UpdateCategoryTemplate)
		admin.DELETE("/:id", controller.DeleteCategoryTemplate)
	}
}
