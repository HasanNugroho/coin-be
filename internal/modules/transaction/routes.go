package transaction

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	r.POST("/income", controller.CreateIncome)
	r.POST("/expense", controller.CreateExpense)
	r.GET("", controller.GetTransactions)
	r.GET("/filter", controller.FilterTransactions)
	r.GET("/:id", controller.GetTransactionByID)
	r.PUT("/:id", controller.UpdateTransaction)
	r.DELETE("/:id", controller.DeleteTransaction)
}
