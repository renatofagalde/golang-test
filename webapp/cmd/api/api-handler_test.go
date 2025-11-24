package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_app_authenticate(t *testing.T) {

	var tests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{name: "valid user", requestBody: `{"username":"admin@example.com","password":"secret"}`, expectedStatusCode: http.StatusCreated},
		{name: "empty user", requestBody: `{"username":"","password":"secret"}`, expectedStatusCode: http.StatusUnauthorized},
		{name: "invalid user", requestBody: `{"username":"nobody@example.com","password":"secret"}`, expectedStatusCode: http.StatusUnauthorized},
		{name: "not json", requestBody: `I'm not a json`, expectedStatusCode: http.StatusUnauthorized},
		{name: "empty json", requestBody: `{}`, expectedStatusCode: http.StatusUnauthorized},
	}

	for _, e := range tests {
		var reader io.Reader
		reader = strings.NewReader(e.requestBody)
		request, _ := http.NewRequest(http.MethodPost, "/auth", reader)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.authenticate)

		handler.ServeHTTP(rr, request)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d, but got %d",
				e.name, e.expectedStatusCode, rr.Code)
		}
	}
}
