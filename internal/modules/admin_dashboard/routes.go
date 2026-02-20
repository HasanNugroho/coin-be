package admin_dashboard

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	r.GET("/dashboard", controller.GetAdminSummary)
}
