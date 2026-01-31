package transaction

import (
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "transactionRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewRepository(client.Database(cfg.MongoDB)), nil
		},
	})

	builder.Add(di.Def{
		Name: "transactionService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("transactionRepository").(*Repository)
			pocketRepo := ctn.Get("pocketRepository").(*pocket.Repository)
			return NewService(repo, pocketRepo), nil
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
