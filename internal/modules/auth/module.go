package auth

import (
	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/category_template"
	"github.com/HasanNugroho/coin-be/internal/modules/platform"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket_template"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/HasanNugroho/coin-be/internal/modules/user_category"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"github.com/redis/go-redis/v9"
	"github.com/sarulabs/di/v2"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "authService",
		Build: func(ctn di.Container) (interface{}, error) {
			userRepo := ctn.Get("userRepository").(*user.Repository)
			pocketRepo := ctn.Get("pocketRepository").(*pocket.Repository)
			pocketTemplateRepo := ctn.Get("pocketTemplateRepository").(*pocket_template.Repository)
			categoryTemplateRepo := ctn.Get("categoryTemplateRepository").(*category_template.Repository)
			userCategoryRepo := ctn.Get("userCategoryRepository").(*user_category.Repository)
			platformRepo := ctn.Get("platformRepository").(*platform.Repository)
			userPlatformRepo := ctn.Get("userPlatformRepository").(*user_platform.UserPlatformRepository)
			redisClient := ctn.Get("redis").(*redis.Client)
			cfg := ctn.Get("config").(*config.Config)
			jwtManager := utils.NewJWTManager(cfg)
			passwordMgr := utils.NewPasswordManager()
			return NewService(userRepo, pocketRepo, pocketTemplateRepo, categoryTemplateRepo, userCategoryRepo, platformRepo, userPlatformRepo, redisClient, jwtManager, passwordMgr), nil
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
