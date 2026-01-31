package container

import (
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/core/database"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
)

func BuildContainer() (di.Container, error) {
	builder, _ := di.NewBuilder()

	// Config
	builder.Add(di.Def{
		Name: "config",
		Build: func(ctn di.Container) (interface{}, error) {
			return config.Load(), nil
		},
	})

	// Mongo
	builder.Add(di.Def{
		Name: "mongo",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			return database.NewMongoClient(cfg.MongoURI)
		},
		Close: func(obj interface{}) error {
			return obj.(*mongo.Client).Disconnect(nil)
		},
	})

	// Redis
	builder.Add(di.Def{
		Name: "redis",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			return database.NewRedisClient(
				cfg.RedisAddr,
				cfg.RedisPass,
				cfg.RedisDB,
			), nil
		},
	})

	// JWT Manager
	builder.Add(di.Def{
		Name: "jwtManager",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			return utils.NewJWTManager(cfg), nil
		},
	})

	// Password Manager
	builder.Add(di.Def{
		Name: "passwordManager",
		Build: func(ctn di.Container) (interface{}, error) {
			return utils.NewPasswordManager(), nil
		},
	})

	return builder.Build(), nil
}
