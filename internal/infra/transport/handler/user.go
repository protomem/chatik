package handler

import (
	"net/http"

	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/internal/core/service"
	"github.com/protomem/chatik/internal/infra/logging"
	"github.com/protomem/chatik/internal/infra/transport"
)

type User struct {
	log  *logging.Logger
	serv service.User
}

func NewUser(log *logging.Logger, serv service.User) *User {
	return &User{
		log:  log.With("handler", "user"),
		serv: serv,
	}
}

func (h *User) Find() http.HandlerFunc {
	type Response struct {
		Users []entity.User `json:"users"`
	}

	return h.makeHandler("find", func(l *logging.Logger, w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		users, err := h.serv.Find(ctx)
		if err != nil {
			return err
		}

		return transport.WriteJSON(w, http.StatusOK, Response{users})
	})
}

func (h *User) makeHandler(endpoint string, apiH APIHandler) http.HandlerFunc {
	return MakeHTTPHandler(h.log.With("endpoint", endpoint), apiH, DefaultErrorHandler)
}
