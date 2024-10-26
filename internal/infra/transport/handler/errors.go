package handler

import (
	"errors"
	"net/http"

	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/internal/infra/logging"
	"github.com/protomem/chatik/internal/infra/transport"
	"github.com/protomem/chatik/pkg/validation"
)

var ErrNotImplemented = errors.New("not implemented")

func DefaultErrorHandler(log *logging.Logger, w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	log.Warn("failed to handle request", logging.Error(err))

	var (
		code = http.StatusInternalServerError
		resp = transport.JSONObject{"error": http.StatusText(code)}
	)

	var verr *validation.Validator
	if errors.As(err, &verr) {
		code = http.StatusUnprocessableEntity
		resp = make(transport.JSONObject, 2)
		if len(verr.Errors) > 0 {
			resp["errors"] = verr.Errors
		}
		if len(verr.FieldErrors) > 0 {
			resp["fieldErrors"] = verr.FieldErrors
		}
	}

	var eerr entity.Error
	if errors.As(err, &eerr) && errors.Is(eerr.Err, entity.ErrNotFound) {
		code = http.StatusNotFound
		resp = transport.JSONObject{"error": eerr.Error()}
	}

	if errors.Is(err, ErrNotImplemented) {
		code = http.StatusNotImplemented
		resp = transport.JSONObject{"error": http.StatusText(http.StatusNotImplemented)}
	}

	if err := transport.WriteJSON(w, code, resp); err != nil {
		log.Error("failed to write response", logging.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
