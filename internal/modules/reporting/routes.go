package reporting

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	// Daily reports
	router.GET("/daily", controller.GetDailyReport)
	router.POST("/daily/generate", controller.GenerateDailyReport)

	// Dashboard endpoints
	dashboard := router.Group("/dashboard")
	{
		dashboard.GET("/kpis", controller.GetRealtimeDashboardKPIs)

		charts := dashboard.Group("/charts")
		{
			charts.GET("/monthly-trend", controller.GetMonthlyTrendChart)
			charts.GET("/pocket-distribution", controller.GetRealtimePocketDistributionChart)
			charts.GET("/expense-by-category", controller.GetRealtimeExpenseByCategoryChart)
		}
	}

	// AI endpoints
	ai := router.Group("/ai")
	{
		ai.GET("/financial-context", controller.GetAIFinancialContext)
	}

	// Health check
	router.GET("/health", controller.HealthCheck)
}
