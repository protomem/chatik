package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/internal/core/service"
	"github.com/protomem/chatik/internal/infra/ctxstore"
	"github.com/protomem/chatik/internal/infra/logging"
	"github.com/protomem/chatik/internal/infra/transport"
)

type Message struct {
	log  *logging.Logger
	serv service.Message
}

func NewMessage(log *logging.Logger, serv service.Message) *Message {
	return &Message{
		log:  log.With("handler", "message"),
		serv: serv,
	}
}

func (h *Message) Find() http.HandlerFunc {
	type Response struct {
		Messages []entity.Message `json:"messages"`
	}

	return h.makeHandler("find", func(log *logging.Logger, w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		fromUser := ctxstore.MustFrom[entity.User](ctx, ctxstore.User)

		toUser, err := entity.ParseID(r.PathValue("userId"))
		if err != nil {
			return transport.WriteError(w, http.StatusBadRequest, errors.New("invalid user id"))
		}

		findOpts, err := h.extractFindOptionsFromRequest(r, service.FindOptions{Limit: 10, Offset: 0})
		if err != nil {
			return transport.WriteError(w, http.StatusBadRequest, err)
		}

		opts := service.FindMessageByFromAndTo{
			FindOptions: findOpts,
			From:        fromUser.ID,
			To:          toUser,
		}

		messages, err := h.serv.FindByFromAndTo(ctx, opts)
		if err != nil {
			return err
		}

		return transport.WriteJSON(w, http.StatusOK, Response{messages})
	})
}

func (h *Message) Create() http.HandlerFunc {
	type Request struct {
		Text string `json:"text"`
	}

	type Response struct {
		entity.Message `json:"message"`
	}

	return h.makeHandler("create", func(l *logging.Logger, w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		fromUser := ctxstore.MustFrom[entity.User](ctx, ctxstore.User)

		toUser, err := entity.ParseID(r.PathValue("userId"))
		if err != nil {
			return transport.WriteError(w, http.StatusBadRequest, errors.New("invalid user id"))
		}

		var req Request
		if err := transport.ReadJSON(w, r, &req); err != nil {
			return err
		}

		dto := service.CreateMessageDTO{
			From: fromUser.ID,
			To:   toUser,
			Text: req.Text,
		}

		message, err := h.serv.Create(ctx, dto)
		if err != nil {
			return err
		}

		return transport.WriteJSON(w, http.StatusCreated, Response{message})
	})
}

func (h *Message) makeHandler(endpoint string, apiH APIHandler) http.HandlerFunc {
	return MakeHTTPHandler(h.log.With("endpoint", endpoint), apiH, DefaultErrorHandler)
}

func (h *Message) extractFindOptionsFromRequest(r *http.Request, defautlValue service.FindOptions) (service.FindOptions, error) {
	var (
		err  error
		opts service.FindOptions = defautlValue
	)

	rawLimit, ok := r.URL.Query().Get("limit"), r.URL.Query().Has("limit")
	if ok {
		opts.Limit, err = strconv.ParseUint(rawLimit, 10, 64)
		if err != nil {
			return opts, errors.New("invalid limit")
		}
	}

	rawOffset, ok := r.URL.Query().Get("offset"), r.URL.Query().Has("offset")
	if ok {
		opts.Offset, err = strconv.ParseUint(rawOffset, 10, 64)
		if err != nil {
			return opts, errors.New("invalid offset")
		}
	}

	return opts, nil
}
