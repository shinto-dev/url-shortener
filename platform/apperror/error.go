package apperror

import (
	"fmt"

	"github.com/pkg/errors"
)

type Code string

type AppError struct {
	Code    Code
	Message string
}

func NewError(code Code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func (a *AppError) Error() string {
	return fmt.Sprintf("error code: %s, message: %s", a.Code, a.Message)
}

func FindError(err error) *AppError {
	var appError *AppError
	errors.As(err, &appError)
	return appError
}
