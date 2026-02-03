package payroll

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
		Name: "payrollRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			return NewRepository(client.Database(cfg.MongoDB)), nil
		},
	})

	builder.Add(di.Def{
		Name: "payrollService",
		Build: func(ctn di.Container) (interface{}, error) {
			payrollRepo := ctn.Get("payrollRepository").(*Repository)
			userRepo := ctn.Get("userRepository").(*user.Repository)
			userPlatformRepo := ctn.Get("userPlatformRepository").(*user_platform.UserPlatformRepository)
			pocketRepo := ctn.Get("pocketRepository").(*pocket.Repository)
			transactionRepo := ctn.Get("transactionRepository").(*transaction.Repository)

			balanceProcessor := transaction.NewBalanceProcessor(pocketRepo, userPlatformRepo)

			cfg := ctn.Get("config").(*config.Config)
			client := ctn.Get("mongo").(*mongo.Client)
			db := client.Database(cfg.MongoDB)
			return NewService(payrollRepo, userRepo, userPlatformRepo, pocketRepo, transactionRepo, balanceProcessor, db), nil
		},
	})
}
