package pocket_template

import (
	"github.com/HasanNugroho/coin-be/internal/core/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	protected := r.Group("")
	{
		protected.GET("", controller.ListPocketTemplates)
		protected.GET("/active", controller.ListActivePocketTemplates)
		protected.GET("/type/:type", controller.ListPocketTemplatesByType)
		protected.GET("/:id", controller.GetPocketTemplate)
	}

	admin := r.Group("")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("", controller.CreatePocketTemplate)
		admin.PUT("/:id", controller.UpdatePocketTemplate)
		admin.DELETE("/:id", controller.DeletePocketTemplate)
	}
}
