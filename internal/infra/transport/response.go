package transport

import (
	"encoding/json"
	"net/http"

	"github.com/protomem/chatik/pkg/werrors"
)

type JSONObject map[string]any

func WriteWithoutBody(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)
	return nil
}

func WriteError(w http.ResponseWriter, status int, err error) error {
	return WriteJSON(w, status, JSONObject{"error": err.Error()})
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	return WriteJSONWithHeaders(w, status, data, nil)
}

func WriteJSONWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
	werr := werrors.Wrap("response/writeJSON")

	encData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return werr(err)
	}

	encData = append(encData, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set(HeaderContentType, MIMETypeJSON)
	w.WriteHeader(status)
	err = werrors.Raise(w.Write(encData))

	return werr(err)
}
