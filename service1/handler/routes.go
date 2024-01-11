package handler

import (
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"

	"service1/pkg/types"
	"service1/pkg/validator"
)

type invoice struct {
	CurrencyCode string  `json:"currencyCode" validate:"required,min=3,max=10"`
	Amount       float64 `json:"amount"       validate:"required,min=0"`
	Wallet       string  `json:"wallet"       validate:"required,min=5,max=50"`
}

func (a *api) invoice(e echo.Context) error {
	ctx := e.Request().Context()
	// Validation and transform type
	valid := validator.New[invoice]()
	input, err := valid.ValidateRequest(ctx, e)
	if err != nil {
		return e.String(http.StatusBadRequest, "err validate")
	}
	transaction := transformInvoiceToTransactions(input)
	// Usecase
	id, status, err := a.app.CreateTransaction(ctx, transaction)
	if err != nil {
		log.Println(err)
		return e.String(http.StatusForbidden, "CreateTransaction failed")
	}

	return e.String(http.StatusOK, fmt.Sprintf("Create invoice transaction %v status %v", *id, *status))
}

type withdraw struct {
	CurrencyCode string  `json:"currencyCode" validate:"required,min=3,max=10"`
	Amount       float64 `json:"amount"       validate:"required,max=0"`
	Wallet       string  `json:"wallet"       validate:"required,min=5,max=50"`
}

func (a *api) withdraw(e echo.Context) error {
	ctx := e.Request().Context()
	// Validation and transform type
	valid := validator.New[withdraw]()
	input, err := valid.ValidateRequest(ctx, e)
	if err != nil {
		return e.String(http.StatusBadRequest, "err validate")
	}
	transaction := transformWithdrawToTransactions(input)
	transaction.IDSender = e.Param("id") // в идеале получаем id из jwt
	// Usecase
	id, status, err := a.app.CreateTransaction(ctx, transaction)
	if err != nil {
		log.Println(err)
		return e.String(http.StatusForbidden, "CreateTransaction failed")
	}

	return e.String(http.StatusOK, fmt.Sprintf("Create withdraw transaction %v status %v", *id, *status))
}

func (a *api) getBalance(e echo.Context) error {
	ctx := e.Request().Context()
	id := e.Param("id") // в идеале получаем id из jwt

	// Usecase
	wallet, err := a.app.GetBalance(ctx, id)
	if err != nil {
		log.Println(err)
		return e.String(http.StatusForbidden, "balance failed")
	}

	return e.String(http.StatusOK, fmt.Sprintf("Balance wallet %v:\n actual balance - %v\n frozen balance - %v\n ",
		wallet.ID, transformIntToFloat(wallet.Actual), transformIntToFloat(wallet.Frozen)))
}

type transactionMessage struct {
	ID           string `json:"id"           validate:"required,min=5,max=50"`
	CurrencyCode string `json:"currencyCode" validate:"required,min=3,max=10"`
	Amount       int    `json:"amount"       validate:"required"`
	Wallet       string `json:"wallet"       validate:"required,min=5,max=50"`
}

func (a *api) processingTransaction(e echo.Context) error {
	ctx := e.Request().Context()
	// Validation and transform type
	valid := validator.New[transactionMessage]()
	input, err := valid.ValidateRequest(ctx, e)
	if err != nil {
		return e.String(http.StatusBadRequest, "err validate")
	}
	transaction := transformTransactionMessageToTransactions(input)
	// Usecase
	if err := a.app.ProcessingTransaction(ctx, transaction); err != nil {
		log.Println(err)
		return e.String(http.StatusForbidden, "Error")
	}

	return e.String(http.StatusOK, "Success")
}

func transformInvoiceToTransactions(transaction *invoice) *types.Transactions {
	return &types.Transactions{
		CurrencyCode: transaction.CurrencyCode,
		Amount:       transformFloatToInt(transaction.Amount),
		Wallet:       transaction.Wallet,
	}
}

func transformWithdrawToTransactions(transaction *withdraw) *types.Transactions {
	return &types.Transactions{
		CurrencyCode: transaction.CurrencyCode,
		Amount:       transformFloatToInt(transaction.Amount),
		Wallet:       transaction.Wallet,
	}
}

func transformTransactionMessageToTransactions(transaction *transactionMessage) *types.Transactions {
	return &types.Transactions{
		ID:           transaction.ID,
		CurrencyCode: transaction.CurrencyCode,
		Amount:       transaction.Amount,
		Wallet:       transaction.Wallet,
	}
}

func transformFloatToInt(a float64) int {
	return int(a * 100)
}

func transformIntToFloat(a int) float64 {
	return float64(a) / 100
}
