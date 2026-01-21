package report

import (
	"github.com/HasanNugroho/coin-be/internal/modules/allocation"
	"github.com/HasanNugroho/coin-be/internal/modules/category"
	"github.com/HasanNugroho/coin-be/internal/modules/target"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"github.com/sarulabs/di/v2"
)

func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "reportService",
		Build: func(ctn di.Container) (interface{}, error) {
			transactionRepo := ctn.Get("transactionRepository").(*transaction.Repository)
			allocationRepo := ctn.Get("allocationRepository").(*allocation.Repository)
			categoryRepo := ctn.Get("categoryRepository").(*category.Repository)
			targetRepo := ctn.Get("targetRepository").(*target.Repository)
			return NewService(transactionRepo, allocationRepo, categoryRepo, targetRepo), nil
		},
	})

	builder.Add(di.Def{
		Name: "reportController",
		Build: func(ctn di.Container) (interface{}, error) {
			svc := ctn.Get("reportService").(*Service)
			return NewController(svc), nil
		},
	})
}
