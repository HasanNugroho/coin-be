package reporting

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all reporting/dashboard routes
func RegisterRoutes(routes *gin.RouterGroup, controller *DashboardController) {
	routes.GET("/kpis", controller.GetKPIs)
	routes.GET("/charts", controller.GetCharts)
	routes.GET("/reports/daily", controller.GetDailyReportsByDateRange)
	routes.GET("/reports/monthly", controller.GetMonthlyReportsByDateRange)
	routes.GET("/charts/income-expense", controller.GetIncomeExpenseChart)
	routes.GET("/charts/category-distribution", controller.GetCategoryDistribution)
	routes.GET("/charts/pocket-distribution", controller.GetPocketDistribution)
	routes.GET("/summary", controller.GetDashboardSummary)
}
