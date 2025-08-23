package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/abc", http.StatusNotFound},
	}

	routes := app.routes()

	//create a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	pathToTemplates = "./../../templates/"

	//range through test data
	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected status %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestAppHome(t *testing.T) {
	//create a request
	request, _ := http.NewRequest("GET", "/", nil)
	request = addContextSessionToRequest(request, &app)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Home)
	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusOK {
		t.Errorf("Test app home return wrong status code, expected 200 but bot %d", rr.Code)
	}

	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), `<small>From Session:`) {
		t.Error("did not find correct text in html")
	}
}

func getCtx(request *http.Request) context.Context {
	ctx := context.WithValue(request.Context(), contextUserKey, "unkwown")
	return ctx
}

func addContextSessionToRequest(request *http.Request, app *application) *http.Request {
	request = request.WithContext(getCtx(request))
	ctx, _ := app.Session.Load(request.Context(), request.Header.Get("X-Session"))
	return request.WithContext(ctx)
}
