package allocation

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	// User routes
	protected := r.Group("")
	{
		protected.POST("", controller.CreateAllocation)
		protected.GET("", controller.ListAllocations)
		protected.GET("/:id", controller.GetAllocation)
		protected.PUT("/:id", controller.UpdateAllocation)
		protected.DELETE("/:id", controller.DeleteAllocation)
	}
}
