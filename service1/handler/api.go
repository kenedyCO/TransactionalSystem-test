package handler

import (
	"context"

	"github.com/labstack/echo"

	"service1/pkg/types"
)

type Application interface {
	GetBalance(ctx context.Context, id string) (*types.Wallet, error)
	CreateTransaction(ctx context.Context, transactions *types.Transactions) (*string, *string, error)
	ProcessingTransaction(ctx context.Context, transactions *types.Transactions) error
}

type api struct {
	app Application
}

func New(usecase Application) *api {
	return &api{
		app: usecase,
	}
}

func (a *api) AddRoute(e *echo.Echo) {
	version := e.Group("/v1")
	{
		service1 := version.Group("/service1")
		{
			service1.POST("/invoice", a.invoice)
			service1.POST("/withdraw/:id", a.withdraw)
			service1.GET("/:id", a.getBalance)
			service1.POST("/transaction", a.processingTransaction)
		}
	}
}
