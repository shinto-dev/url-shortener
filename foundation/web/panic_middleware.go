package web

import (
	"net/http"

	"github.com/pkg/errors"
)

func PanicHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				HandleError(r.Context(), w, getPanicErr(r))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
func getPanicErr(recoverErr interface{}) error {
	err, ok := recoverErr.(error)
	if !ok {
		return errors.Errorf("%v", recoverErr)
	}
	return err
}
