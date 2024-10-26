package handler

import (
	"net/http"

	"github.com/protomem/chatik/internal/infra/transport"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	transport.WriteJSON(w, http.StatusOK, transport.JSONObject{"status": "OK"})
}
