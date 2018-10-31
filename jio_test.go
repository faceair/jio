package jio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateJSON(t *testing.T) {
	data := []byte(`{`)
	if _, err := ValidateJSON(&data, Any()); err == nil {
		t.Error("should error")
	}

	data = []byte(`{"1": 10}`)
	if _, err := ValidateJSON(&data, Object().Keys(K{"1": Number().Max(5)})); err == nil {
		t.Error("should error")
	}

	data = []byte(`{"1": 10}`)
	if _, err := ValidateJSON(&data, Object().Keys(K{"1": Any().Transform(func(ctx *Context) {
		ctx.Value = make(chan int)
	})})); err == nil {
		t.Error("should error")
	}

	data = []byte(`{"1": 10}`)
	if _, err := ValidateJSON(&data, Object().Keys(K{"1": Number()})); err != nil {
		t.Error("should no error")
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestValidateBody(t *testing.T) {
	schema := Object().Keys(K{
		"debug": Bool().Truthy("on").Required(),
		"window": Object().Keys(K{
			"title": String().Min(3).Max(18).Required(),
			"size":  Array().Items(Number().Integer()).Length(2).Required(),
		}).Without("name", "title").Required(),
	})
	handler := ValidateBody(schema, DefaultErrorHandler)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Post(ts.URL, "application/json", strings.NewReader(`{
		"debug": "on",
		"window": {
			"title": "Sample Widget",
			"size": [500, 500]
		}
	}`))
	if err != nil {
		t.Error(err.Error())
	}
	if res != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	if string(body) != "ok" {
		t.Error("not ok")
	}

	res, err = http.Post(ts.URL, "application/json", strings.NewReader(`{
		"debug": "on",
		"window": {
			"title": "Sample Widget",
			"size": [500]
		}
	}`))
	if err != nil {
		t.Error(err.Error())
	}
	if res != nil {
		defer res.Body.Close()
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Error("should bad request")
	}

	testRequest := httptest.NewRequest(http.MethodPost, "/something", errReader(0))
	testRequest.Header.Add("Content-Type", "application/json")
	handler.ServeHTTP(nil, testRequest)
}

func TestValidateQuery(t *testing.T) {
	schema := Object().Keys(K{
		"keyword": String(),
		"limit":   Number().ParseString().Integer(),
	})
	handler := ValidateQuery(schema, DefaultErrorHandler)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Post(ts.URL+"?keyword=test&limit=1", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Error(err.Error())
	}
	if res != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
	}
	if string(body) != "ok" {
		t.Error("not ok")
	}

	res, err = http.Post(ts.URL+"?keyword=test&limit=1.1", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Error(err.Error())
	}
	if res != nil {
		defer res.Body.Close()
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Error("should bad request")
	}
}
