package handlers

import (
	"context"
	"net/http"

	"github.com/shinto-dev/url-shortener/foundation/apperror"
	"github.com/shinto-dev/url-shortener/foundation/observation/logging"
	"github.com/shinto-dev/url-shortener/foundation/web"
	"go.uber.org/zap"
)

const (
	AppErrInvalidRequestBody = "invalid_request_body"
)

const (
	ErrInvalidInput         = "ERR-101"
	ErrCodeShortURLNotFound = "ERR-201"
)

var (
	ErrShortURLNotFound = web.NewRequestError(http.StatusNotFound, web.Error{
		Code:    ErrCodeShortURLNotFound,
		Message: "short url not found",
	})
	ErrInvalidRequestBody = web.NewRequestError(http.StatusBadRequest, web.Error{
		Code:    ErrInvalidInput,
		Message: "request body is invalid",
	})
)

func commonErrHandler(errMap map[apperror.Code]*web.RequestError) func(ctx context.Context, err error) error {
	return func(ctx context.Context, err error) error {
		appErr := apperror.FindError(err)
		if appErr == nil {
			logging.FromContext(ctx).Error("unexpected error!", err)
			return err
		}

		logging.FromContext(ctx).
			WithFields(zap.Error(appErr)).
			Info("error while handling request")

		if err, ok := errMap[appErr.Code]; ok {
			return err
		}

		logging.FromContext(ctx).
			WithFields(zap.Any("code", appErr.Code)).
			Warn("unmapped error code")

		return err
	}
}

func validationErrorHandler(ctx context.Context, err error) error {
	logging.FromContext(ctx).
		WithFields(zap.Error(err)).Info("validation failed")

	return web.NewRequestError(http.StatusBadRequest, web.Error{
		Code:    ErrInvalidInput,
		Message: err.Error(),
	})
}
