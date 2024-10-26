package handler

import (
	"net/http"

	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/internal/core/service"
	"github.com/protomem/chatik/internal/infra/ctxstore"
	"github.com/protomem/chatik/internal/infra/logging"
	"github.com/protomem/chatik/internal/infra/transport"
)

type Auth struct {
	log  *logging.Logger
	serv service.Auth
}

func NewAuth(log *logging.Logger, serv service.Auth) *Auth {
	return &Auth{
		log:  log.With("handler", "auth"),
		serv: serv,
	}
}

func (h *Auth) Login() http.HandlerFunc {
	type Request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	type Response struct {
		entity.User       `json:"user"`
		service.TokensDTO `json:",inline"`
	}

	return h.makeHandler("login", func(log *logging.Logger, w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		var req Request
		if err := transport.ReadJSON(w, r, &req); err != nil {
			return err
		}

		dto := service.LoginDTO{
			Login:       req.Login,
			Password:    req.Password,
			Fingerprint: ctxstore.MustFrom[string](ctx, ctxstore.ClientIP),
		}

		user, tokens, err := h.serv.Login(ctx, dto)
		if err != nil {
			return err
		}

		return transport.WriteJSON(w, http.StatusOK, Response{user, tokens})
	})
}

func (h *Auth) Register() http.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type Response struct {
		entity.User       `json:"user"`
		service.TokensDTO `json:",inline"`
	}

	return h.makeHandler("register", func(log *logging.Logger, w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		var req Request
		if err := transport.ReadJSON(w, r, &req); err != nil {
			return err
		}

		dto := service.RegisterDTO{
			CreateUserDTO: service.CreateUserDTO(req),
			Fingerprint:   ctxstore.MustFrom[string](ctx, ctxstore.ClientIP),
		}

		user, tokens, err := h.serv.Register(ctx, dto)
		if err != nil {
			return err
		}

		return transport.WriteJSON(w, http.StatusCreated, Response{user, tokens})
	})
}

func (h *Auth) Logout() http.HandlerFunc {
	return h.makeHandler("logout", func(log *logging.Logger, w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		sessionToken := r.Header.Get(transport.HeaderXSessionToken)

		if err := h.serv.Logout(ctx, sessionToken); err != nil {
			return err
		}

		return transport.WriteWithoutBody(w, http.StatusNoContent)
	})
}

func (h *Auth) Refresh() http.HandlerFunc {
	type Response struct {
		service.TokensDTO `json:",inline"`
	}

	return h.makeHandler("refresh", func(log *logging.Logger, w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		sessionToken := r.Header.Get(transport.HeaderXSessionToken)

		dto := service.RefreshDTO{
			Token:       sessionToken,
			Fingerprint: ctxstore.MustFrom[string](ctx, ctxstore.ClientIP),
		}

		tokens, err := h.serv.Refresh(ctx, dto)
		if err != nil {
			return err
		}

		return transport.WriteJSON(w, http.StatusOK, Response{tokens})
	})
}

func (h *Auth) makeHandler(endpoint string, apiH APIHandler) http.HandlerFunc {
	return MakeHTTPHandler(h.log.With("endpoint", endpoint), apiH, DefaultErrorHandler)
}
