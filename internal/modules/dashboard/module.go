package dashboard

import (
	"context"

	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/modules/daily_summary"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "dashboardRepository",
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
		Name: "dashboardService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("dashboardRepository").(*Repository)
			dsr := ctn.Get("dailySummaryRepository").(*daily_summary.Repository)
			return NewService(repo, dsr), nil
		},
	})

	builder.Add(di.Def{
		Name: "dashboardController",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("dashboardService").(*Service)
			dss := ctn.Get("dailySummaryService").(*daily_summary.Service)
			return NewController(service, dss), nil
		},
	})
}
