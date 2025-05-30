package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
)

type envelope map[string]interface{}

func (app *application) writeJSON(
	w http.ResponseWriter,
	status int,
	data envelope,
	headers http.Header,
) error {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	json = append(json, '\n')

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}

func (app *application) readJSON(
	w http.ResponseWriter,
	r *http.Request,
	destination interface{},
) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(destination)
	if err != nil {
		var (
			syntaxError           *json.SyntaxError
			unmarshalTypeError    *json.UnmarshalTypeError
			invalidUnmarshalError *json.InvalidUnmarshalError
		)
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf(
				"body contains badly-formed JSON (at character %d)",
				syntaxError.Offset,
			)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf(
					"body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field,
				)
			}
			return fmt.Errorf(
				"body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset,
			)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: uknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func removeFromSliceByValue[T comparable](s []T, val T) []T {
	for i, v := range s {
		if v == val {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s // not found; return unchanged
}
