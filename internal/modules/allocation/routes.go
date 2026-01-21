package allocation

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	r.POST("", controller.CreateAllocation)
	r.GET("", controller.GetAllocations)
	r.GET("/:id", controller.GetAllocationByID)
	r.PUT("/:id", controller.UpdateAllocation)
	r.DELETE("/:id", controller.DeleteAllocation)
	r.GET("/logs", controller.GetAllocationLogs)
	r.GET("/:id/logs", controller.GetAllocationLogsByAllocationID)
}
