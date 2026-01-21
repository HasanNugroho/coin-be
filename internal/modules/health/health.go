package health

import (
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Check(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "healthController",
		Build: func(ctn di.Container) (interface{}, error) {
			return NewController(), nil
		},
	})
}

func RegisterRoutes(r *gin.RouterGroup, controller *Controller) {
	r.GET("/", controller.Check)
}
