package dashboard

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	protected := r.Group("")
	{
		protected.GET("/summary", controller.GetDashboardSummary)
		protected.GET("/charts", controller.GetDashboardCharts)
	}
}
