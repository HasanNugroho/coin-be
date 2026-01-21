package target

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	r.POST("", controller.CreateTarget)
	r.GET("", controller.GetTargets)
	r.GET("/:id", controller.GetTargetByID)
	r.PUT("/:id", controller.UpdateTarget)
	r.DELETE("/:id", controller.DeleteTarget)
}
