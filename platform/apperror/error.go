package apperror

import (
	"fmt"

	"github.com/pkg/errors"
)

type Code string

type AppError struct {
	Code    Code
	Message string
	Cause   error
}

func NewError(code Code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewErrorWithCause(err error, code Code, message string) *AppError {
	return &AppError{Code: code, Message: message, Cause: err}
}

func (a *AppError) Error() string {
	if a.Cause != nil {
		return fmt.Sprintf("error code: %s, message: %s, cause: %s", a.Code, a.Message, a.Cause.Error())
	}

	return fmt.Sprintf("error code: %s, message: %s", a.Code, a.Message)
}

func FindError(err error) *AppError {
	var appError *AppError
	errors.As(err, &appError)
	return appError
}

func Is(err error, code Code) bool {
	appErr := FindError(err)
	if appErr == nil {
		return false
	}

	return appErr.Code == code
}
