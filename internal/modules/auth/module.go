package auth

import (
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/redis/go-redis/v9"
	"github.com/sarulabs/di/v2"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "authService",
		Build: func(ctn di.Container) (interface{}, error) {
			userRepo := ctn.Get("userRepository").(*user.Repository)
			redisClient := ctn.Get("redis").(*redis.Client)
			cfg := ctn.Get("config").(*config.Config)
			jwtManager := utils.NewJWTManager(cfg.JWTSecret)
			passwordMgr := utils.NewPasswordManager()
			return NewService(userRepo, redisClient, jwtManager, passwordMgr), nil
		},
	})

	builder.Add(di.Def{
		Name: "authController",
		Build: func(ctn di.Container) (interface{}, error) {
			svc := ctn.Get("authService").(*Service)
			userSrv := ctn.Get("userService").(*user.Service)
			return NewController(svc, userSrv), nil
		},
	})
}
