package transaction

import (
	"context"

	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/modules/dashboard"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "transactionRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			repo := NewRepository(client.Database(cfg.MongoDB))

			if err := repo.EnsureIndexes(context.Background()); err != nil {
				return nil, err
			}

			return repo, nil
		},
	})

	builder.Add(di.Def{
		Name: "transactionService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("transactionRepository").(*Repository)
			pocketRepo := ctn.Get("pocketRepository").(*pocket.Repository)
			userPlatformRepo := ctn.Get("userPlatformRepository").(*user_platform.UserPlatformRepository)
			dashboardService := ctn.Get("dashboardService").(*dashboard.Service)
			return NewService(repo, pocketRepo, userPlatformRepo, dashboardService), nil
		},
	})

	builder.Add(di.Def{
		Name: "transactionController",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("transactionService").(*Service)
			return NewController(service), nil
		},
	})
}
