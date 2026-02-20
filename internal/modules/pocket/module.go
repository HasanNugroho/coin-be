package pocket

import (
	"context"

	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "pocketRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			repo := NewRepository(client.Database(cfg.MongoDB))
			repo.EnsureIndexes(context.Background())

			return repo, nil
		},
	})

	builder.Add(di.Def{
		Name: "pocketService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("pocketRepository").(*Repository)
			return NewService(repo), nil
		},
	})

	builder.Add(di.Def{
		Name: "pocketController",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("pocketService").(*Service)
			return NewController(service), nil
		},
	})
}
