package reporting

import (
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HasanNugroho/coin-be/internal/core/config"
)

func Register(builder *di.Builder) error {
	builder.Add(di.Def{
		Name: "reportingService",
		Build: func(ctn di.Container) (interface{}, error) {
			mongoClient := ctn.Get("mongo").(*mongo.Client)
			cfg := ctn.Get("config").(*config.Config)
			db := mongoClient.Database(cfg.MongoDB)
			return NewService(db), nil
		},
	})

	builder.Add(di.Def{
		Name: "reportingController",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("reportingService").(*Service)
			return NewController(service), nil
		},
	})

	return nil
}
