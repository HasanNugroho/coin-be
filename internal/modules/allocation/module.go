package allocation

import (
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
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
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			db := client.Database(cfg.MongoDB)
			repo := ctn.Get("allocationRepository").(*Repository)
			pocketRepo := ctn.Get("pocketRepository").(*pocket.Repository)
			userPlatformRepo := ctn.Get("userPlatformRepository").(*user_platform.UserPlatformRepository)
			userRepo := ctn.Get("userRepository").(*user.Repository)
			transactionRepo := ctn.Get("transactionRepository").(*transaction.Repository)
			return NewService(repo, pocketRepo, userPlatformRepo, userRepo, transactionRepo, db), nil
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
