package web

import (
	"context"
	"encoding/json"
	"net/http"
	"url-shortner/platform/apperror"
	"url-shortner/platform/observation/logging"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
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

	_ = JSON(w, http.StatusBadRequest, Error{
		Code:    string(appErr.Code),
		Message: appErr.Message,
	})
}
