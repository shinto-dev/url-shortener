package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gorilla/mux"
)

func TestPanicHandlerMiddleware(t *testing.T) {
	tests := []struct {
		name               string
		HandlerFunc        http.HandlerFunc
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name: "should handle panic",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic("panic")
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"code":    "",
				"message": "internal error",
			},
		},
		{
			name: "should handle normal requests successfully",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				_ = JSON(w, http.StatusOK, map[string]string{
					"message": "Hello World",
				})
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "Hello World",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := mux.NewRouter()
			r.Use(PanicHandlerMiddleware)
			r.HandleFunc("/op", test.HandlerFunc)

			e := httpexpect.New(t, httptest.NewServer(r).URL)

			e.GET("/op").
				Expect().
				Status(test.expectedStatusCode).
				JSON().Equal(test.expectedResponse)
		})
	}
}
