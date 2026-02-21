package transaction

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	protected := r.Group("")
	{
		protected.POST("", controller.CreateTransaction)
		protected.GET("", controller.ListUserTransactions)
		protected.GET("/:id", controller.GetTransaction)
		protected.PUT("/:id", controller.UpdateTransaction)
		protected.DELETE("/:id", controller.DeleteTransaction)
		protected.GET("/pocket/:pocket_id", controller.ListPocketTransactions)
	}
}
