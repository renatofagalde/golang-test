package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestAppHome_V2(t *testing.T) {

	var tests = []struct {
		name         string
		putInSection string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session:"},
	}

	for _, e := range tests {
		request, _ := http.NewRequest("GET", "/", nil)
		request = addContextSessionToRequest(request, &app)

		_ = app.Session.Destroy(request.Context())
		if e.putInSection != "" {
			app.Session.Put(request.Context(), "test", e.putInSection)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Home)

		handler.ServeHTTP(rr, request)

		if rr.Code != http.StatusOK {
			t.Errorf("TestAppHome returned wrong status code; expected 200 but got %d", rr.Code)
		}

		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), e.expectedHTML) {
			t.Errorf("%s: did not find %s, in response body", e.name, e.expectedHTML)
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

func TestApp_renderWithBadTemplate(t *testing.T) {
	//set templatepath to a locatioion with a bad template ft

	pathToTemplates = "./testdata/"
	request, _ := http.NewRequest("GET", "/", nil)

	request = addContextSessionToRequest(request, &app)
	rr := httptest.NewRecorder()

	err := app.render(rr, request, "bad.page.gohtml", &TemplateData{})
	if err == nil {
		t.Error("expected error from bad template")
	}

	pathToTemplates = "./../../templates/"

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

func Test_app_Login(t *testing.T) {
	var tests = []struct {
		name               string
		postedData         url.Values
		expectedStatusCode int
		expectedLoc        string
	}{
		{name: "valid login",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/u/p",
		},
	}

	for _, e := range tests {
		request, _ := http.NewRequest("POST", "/login", strings.NewReader(e.postedData.Encode()))

		request = addContextSessionToRequest(request, &app)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Login)

		handler.ServeHTTP(rr, request)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: returned wrong status code; expected %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}
