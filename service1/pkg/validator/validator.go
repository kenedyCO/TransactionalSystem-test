package validator

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type Validator[T any] struct {
}

func New[T any]() *Validator[T] {
	return &Validator[T]{}
}

func (*Validator[T]) ValidateRequest(ctx context.Context, e echo.Context) (*T, error) {
	var input *T
	if err := e.Bind(&input); err != nil {
		return nil, fmt.Errorf("invalid body: %w", err)
	}

	validate := validator.New()
	if err := validate.StructCtx(ctx, input); err != nil {
		err := err.(validator.ValidationErrors)
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return input, nil
}
