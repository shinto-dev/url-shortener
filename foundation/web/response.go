package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	mediaTypeApplicationJson = "application/json"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type RequestError struct {
	Err    Error
	Status int
}

func NewRequestError(status int, err Error) *RequestError {
	return &RequestError{
		Err:    err,
		Status: status,
	}
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("code: %s, message: %s, status:%d", e.Err.Code, e.Err.Message, e.Status)
}

func findError(err error) *RequestError {
	var appError *RequestError
	errors.As(err, &appError)
	return appError
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

func HandleError(_ context.Context, w http.ResponseWriter, err error) {
	requestError := findError(err)
	if requestError == nil {
		_ = JSON(w, http.StatusInternalServerError, Error{
			Message: "internal error",
		})
		return
	}

	_ = JSON(w, requestError.Status, requestError.Err)
}
