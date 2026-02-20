package admin_dashboard

import (
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "adminDashboardRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewRepository(client.Database(cfg.MongoDB)), nil
		},
	})

	builder.Add(di.Def{
		Name: "adminDashboardService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("adminDashboardRepository").(*Repository)
			return NewService(repo), nil
		},
	})

	builder.Add(di.Def{
		Name: "adminDashboardController",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("adminDashboardService").(*Service)
			return NewController(service), nil
		},
	})
}
