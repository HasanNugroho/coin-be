package target

import (
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/modules/allocation"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "targetRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewRepository(client.Database(cfg.MongoDB)), nil
		},
	})

	builder.Add(di.Def{
		Name: "targetService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("targetRepository").(*Repository)
			allocationRepo := ctn.Get("allocationRepository").(*allocation.Repository)
			return NewService(repo, allocationRepo), nil
		},
	})

	builder.Add(di.Def{
		Name: "targetController",
		Build: func(ctn di.Container) (interface{}, error) {
			svc := ctn.Get("targetService").(*Service)
			return NewController(svc), nil
		},
	})
}
