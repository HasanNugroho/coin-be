package pocket

import (
	"github.com/HasanNugroho/coin-be/internal/core/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	// User routes
	protected := r.Group("")
	{
		protected.POST("", controller.CreatePocket)
		protected.GET("", controller.ListPockets)
		protected.GET("/main", controller.GetMainPocket)
		protected.GET("/active", controller.ListActivePockets)
		protected.GET("/:id", controller.GetPocket)
		protected.PUT("/:id", controller.UpdatePocket)
		protected.DELETE("/:id", controller.DeletePocket)
	}

	// Admin routes
	admin := r.Group("admin")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("/:user_id", controller.CreateSystemPocket)
		admin.GET("", controller.ListAllPockets)
	}
}
