package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/protomem/chatik/pkg/werrors"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return readJSON(w, r, dst, false)
}

func ReadJSONStrict(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return readJSON(w, r, dst, true)
}

func readJSON(w http.ResponseWriter, r *http.Request, dst interface{}, disallowUnknownFields bool) error {
	werr := werrors.Wrap("request/readJSON")

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	if disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return werr(fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset))

		case errors.Is(err, io.ErrUnexpectedEOF):
			return werr(errors.New("body contains badly-formed JSON"))

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return werr(fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field))
			}
			return werr(fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset))

		case errors.Is(err, io.EOF):
			return werr(errors.New("body must not be empty"))

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return werr(fmt.Errorf("body contains unknown key %s", fieldName))

		case err.Error() == "http: request body too large":
			return werr(fmt.Errorf("body must not be larger than %d bytes", maxBytes))

		case errors.As(err, &invalidUnmarshalError):
			panic(werr(err))

		default:
			return werr(err)
		}
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return werr(errors.New("body must only contain a single JSON value"))
	}

	return nil
}
