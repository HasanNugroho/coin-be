package report

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	r.GET("/dashboard", controller.GetDashboardSummary)
	r.GET("/income", controller.GetIncomeReport)
	r.GET("/expense", controller.GetExpenseReport)
	r.GET("/allocation", controller.GetAllocationReport)
	r.GET("/target-progress", controller.GetTargetProgress)
}
