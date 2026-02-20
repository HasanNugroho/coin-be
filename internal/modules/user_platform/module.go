package user_platform

import (
	"context"

	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/modules/platform"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "userPlatformRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			repo := NewUserPlatformRepository(client.Database(cfg.MongoDB))
			repo.EnsureIndexes(context.Background())

			return repo, nil
		},
	})

	builder.Add(di.Def{
		Name: "userPlatformService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("userPlatformRepository").(*UserPlatformRepository)
			platformRepo := ctn.Get("platformRepository").(*platform.Repository)
			return NewService(repo, platformRepo), nil
		},
	})

	builder.Add(di.Def{
		Name: "userPlatformController",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("userPlatformService").(*Service)
			platformRepo := ctn.Get("platformRepository").(*platform.Repository)
			return NewController(service, platformRepo), nil
		},
	})
}
