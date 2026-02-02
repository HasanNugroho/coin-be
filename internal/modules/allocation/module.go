package allocation

import (
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "allocationRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewRepository(client.Database(cfg.MongoDB)), nil
		},
	})

	builder.Add(di.Def{
		Name: "allocationService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("allocationRepository").(*Repository)
			pocketRepo := ctn.Get("pocketRepository").(*pocket.Repository)
			userPlatformRepo := ctn.Get("userPlatformRepository").(*user_platform.UserPlatformRepository)
			return NewService(repo, pocketRepo, userPlatformRepo), nil
		},
	})

	builder.Add(di.Def{
		Name: "allocationController",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("allocationService").(*Service)
			return NewController(service), nil
		},
	})
}
