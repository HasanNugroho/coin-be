package reporting

import (
	"context"

	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "reportingRepository",
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
		Name: "reportingAggregationHelper",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewAggregationHelper(client.Database(cfg.MongoDB)), nil
		},
	})
}

// EnsureIndexes creates all necessary indexes for reporting collections
func (r *Repository) EnsureIndexes(ctx context.Context) error {
	indexManager := NewIndexManager(r.db)
	return indexManager.CreateAllIndexes(ctx)
}
