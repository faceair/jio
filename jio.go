package jio

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type ContextKey int

const (
	ContextKeyQuery ContextKey = iota
	ContextKeyBody
)

func ValidateJSON(dataRaw *[]byte, schema Schema) (dataMap map[string]interface{}, err error) {
	if err = json.Unmarshal(*dataRaw, &dataMap); err != nil {
		return
	}
	ctx := NewContext(dataMap)
	schema.Validate(ctx)
	if ctx.Err != nil {
		return dataMap, ctx.Err
	}
	dataMap = ctx.Value.(map[string]interface{})
	dataNew, err := json.Marshal(ctx.Value)
	if err != nil {
		return
	}
	*dataRaw = dataNew
	return
}

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	code := http.StatusBadRequest
	body, _ := json.Marshal(map[string]string{
		"message": err.Error(),
	})
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	w.Write(body)
}

func ValidateBody(schema Schema, errorHandler func(http.ResponseWriter, *http.Request, error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var body []byte
			var err error
			if strings.Contains(r.Header.Get("Content-type"), "application/json") {
				body, err = ioutil.ReadAll(r.Body)
				if err != nil {
					return
				}
				r.Body.Close()
			}
			dataMap, err := ValidateJSON(&body, schema)
			if err != nil {
				errorHandler(w, r, err)
				return
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ContextKeyBody, dataMap)))
		}
		return http.HandlerFunc(fn)
	}
}

func ValidateQuery(schema Schema, errorHandler func(http.ResponseWriter, *http.Request, error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			query := make(map[string]interface{})
			for key, value := range r.URL.Query() {
				query[key] = value[0]
			}
			ctx := NewContext(query)
			schema.Validate(ctx)
			if ctx.Err != nil {
				errorHandler(w, r, ctx.Err)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ContextKeyQuery, ctx.Value)))
		}
		return http.HandlerFunc(fn)
	}
}
