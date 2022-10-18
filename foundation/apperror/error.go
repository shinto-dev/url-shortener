package apperror

import (
	"fmt"

	"github.com/pkg/errors"
)

type Code string

type AppError struct {
	Code  Code
	Cause error
}

func NewError(code Code) *AppError {
	return &AppError{Code: code}
}

func NewErrorWithCause(err error, code Code) *AppError {
	return &AppError{Code: code, Cause: err}
}

func (a *AppError) Error() string {
	if a.Cause != nil {
		return fmt.Sprintf("error code: %s, cause: %s", a.Code, a.Cause.Error())
	}

	return fmt.Sprintf("error code: %s", a.Code)
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
