package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"cafe_store.hiyabnako/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// envelope map to wrap json responses
type envelope map[string]any

// extract id from path path parameter if valid
func (app *application) readIDparam(r *http.Request) (int64, error) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("book_id"), 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {

	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	jsDec := json.NewDecoder(r.Body)
	jsDec.DisallowUnknownFields()
	
	err := jsDec.Decode(dest)
	if err != nil {

		var (
			syntaxError *json.SyntaxError
			InvalidUnmarshalError *json.InvalidUnmarshalError
			unmarshalTypeError *json.UnmarshalTypeError
			maxBytesError *http.MaxBytesError
		)

		switch {
		
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("the body exceeds the maximum allowed bytes %d", maxBytesError.Limit)
		
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("Body contains unknown key %s", fieldName)

		case errors.As(err, &syntaxError):
			return fmt.Errorf("There is an issue at character %d", syntaxError.Offset)
		
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains a badly formed JSON")

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &InvalidUnmarshalError):
			panic(err)
		
		default:
			return err
		}
    }


	err = jsDec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("Body must only contain one JSON value")
	}

	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, data envelope, headers http.Header) error {

	js := json.NewEncoder(w)

	err := js.Encode(data)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "sorry, server could not process your request", http.StatusInternalServerError)
	}

	
	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")

	return nil
}

func (app *application) readString(qs url.Values, Key string, defaultValue string) string {

	s := qs.Get(Key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, Key string, defaultValue []string) []string {

	s := qs.Get(Key)

	if s == "" {
		return defaultValue
	}

	return strings.Split(s, ",")
}


func (app *application) readInt(qs url.Values, Key string, defaultValue int, v *validator.Validator) int{

	s := qs.Get(Key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddErrors(Key, "must be an integer")
		return defaultValue
	}

	return i
}