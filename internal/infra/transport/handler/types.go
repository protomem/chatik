package handler

import (
	"net/http"

	"github.com/protomem/chatik/internal/infra/logging"
)

type APIHandler func(*logging.Logger, http.ResponseWriter, *http.Request) error

type ErrorHandler func(*logging.Logger, http.ResponseWriter, *http.Request, error)

func MakeHTTPHandler(log *logging.Logger, apiH APIHandler, errH ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := apiH(log.WithContext(ctx), w, r); err != nil {
			errH(log.WithContext(ctx), w, r, err)
		}
	}
}
