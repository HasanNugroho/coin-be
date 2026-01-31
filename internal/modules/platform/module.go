package platform

import (
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "platformRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewRepository(client.Database(cfg.MongoDB)), nil
		},
	})

	builder.Add(di.Def{
		Name: "platformService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("platformRepository").(*Repository)
			return NewService(repo), nil
		},
	})

	builder.Add(di.Def{
		Name: "platformController",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("platformService").(*Service)
			return NewController(service), nil
		},
	})
}
