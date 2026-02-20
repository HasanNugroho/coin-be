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
		protected.GET("/dropdown", controller.ListPocketsDropdown)
		protected.GET("/:id", controller.GetPocket)
		protected.PUT("/:id", controller.UpdatePocket)
		protected.PUT("/:id/lock", controller.LockPocket)
		protected.PUT("/:id/unlock", controller.UnlockPocket)
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
