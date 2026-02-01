package user_category

import (
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HasanNugroho/coin-be/internal/core/config"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "userCategoryRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewRepository(client.Database(cfg.MongoDB)), nil
		},
	})

	builder.Add(di.Def{
		Name: "userCategoryService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("userCategoryRepository").(*Repository)
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewService(repo, client.Database(cfg.MongoDB)), nil
		},
	})

	builder.Add(di.Def{
		Name: "userCategoryController",
		Build: func(ctn di.Container) (interface{}, error) {
			svc := ctn.Get("userCategoryService").(*Service)
			return NewController(svc), nil
		},
	})
}
