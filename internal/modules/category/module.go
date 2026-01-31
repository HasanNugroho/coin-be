package category

import (
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HasanNugroho/coin-be/internal/core/config"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "categoryRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewRepository(client.Database(cfg.MongoDB)), nil
		},
	})

	builder.Add(di.Def{
		Name: "categoryService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("categoryRepository").(*Repository)
			return NewService(repo), nil
		},
	})

	builder.Add(di.Def{
		Name: "categoryController",
		Build: func(ctn di.Container) (interface{}, error) {
			svc := ctn.Get("categoryService").(*Service)
			return NewController(svc), nil
		},
	})
}
