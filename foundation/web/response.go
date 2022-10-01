package web

import (
	"context"
	"encoding/json"
	"net/http"
	"url-shortener/foundation/apperror"
	"url-shortener/foundation/observation/logging"

	"go.uber.org/zap"
)

const (
	mediaTypeApplicationJson = "application/json"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", mediaTypeApplicationJson)
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func Status(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)
	return nil
}

func HandleError(ctx context.Context, w http.ResponseWriter, err error) {
	appErr := apperror.FindError(err)
	if appErr == nil {
		logging.FromContext(ctx).Error("unexpected error!", err)
		_ = JSON(w, http.StatusInternalServerError, Error{
			Message: "internal error",
		})
		return
	}

	logging.FromContext(ctx).
		WithFields(zap.Error(err)).
		Error("error while handling request")

	switch appErr.Code {
	case ErrCodeRecordNotFound:
		_ = Status(w, http.StatusNotFound)
	default:
		_ = JSON(w, http.StatusBadRequest, Error{
			Code:    string(appErr.Code),
			Message: appErr.Message,
		})
	}
}
