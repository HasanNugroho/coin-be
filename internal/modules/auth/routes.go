package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.POST("/refresh", controller.RefreshToken)
	r.POST("/logout", controller.Logout)
}
