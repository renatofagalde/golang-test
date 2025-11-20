package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name                    string
		url                     string
		expectedStatusCode      int
		expectedURL             string
		expectedFirstStatusCode int
	}{
		{"home", "/", http.StatusOK, "/", http.StatusOK},
		{"404", "/abc", http.StatusNotFound, "/fish", http.StatusNotFound},
		{"profile", "/u/p", http.StatusOK, "/", http.StatusTemporaryRedirect},
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

func TestAppHome_V2(t *testing.T) {

	var tests = []struct {
		name         string
		putInSection string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session:"},
		{"second", "hello world!", "<small>"},
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

func getCtx(request *http.Request) context.Context {
	ctx := context.WithValue(request.Context(), contextUserKey, "unkwown")
	return ctx
}

func addContextSessionToRequest(request *http.Request, app *application) *http.Request {
	request = request.WithContext(getCtx(request))
	ctx, _ := app.Session.Load(request.Context(), request.Header.Get("X-Session"))
	return request.WithContext(ctx)
}

func Test_app_UploadFile(t *testing.T) {
	//setup pipes
	pr, pw := io.Pipe()

	//create new writer, of type *io.Writer
	writer := multipart.NewWriter(pw)

	//create a waitgroup, and add 1 to it
	wg := &sync.WaitGroup{}
	wg.Add(1)

	//simulate uploading a file using a goroutine and our writer
	go simulatePingUpload("./testdata/img.png", writer, t, wg)

	//read from the pipe which receives data
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	//call app.UploadFiles
	uploadFiles, err := app.UploadFile(request, "./testdata/uploads/")
	if err != nil {
		t.Error(err)
	}

	//perform our tests
	if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadFiles[0].OriginalFileName)); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", err.Error())
	}

	//clean up
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadFiles[0].OriginalFileName))

	wg.Wait()
}

func simulatePingUpload(fileToUpload string, writer *multipart.Writer, t *testing.T, wg *sync.WaitGroup) {
	defer writer.Close()
	defer wg.Done()

	// 1) Cria um arquivo no disco
	f, err := os.Create(fileToUpload)
	if err != nil {
		t.Fatalf("erro criando arquivo temporário: %v", err)
	}

	// 2) Gera uma imagem PNG simples
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// exemplo preenchendo de azul
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}

	// 3) Grava o PNG dentro do arquivo
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("erro salvando PNG: %v", err)
	}
	f.Close()

	// 4) Cria o multipart field
	part, err := writer.CreateFormFile("file", path.Base(fileToUpload))
	if err != nil {
		t.Fatalf("erro criando form file: %v", err)
	}

	// 5) Reabre o arquivo criado
	f2, err := os.Open(fileToUpload)
	if err != nil {
		t.Fatalf("erro reabrindo arquivo: %v", err)
	}
	defer f2.Close()

	// 6) Copia o conteúdo do PNG para o multipart
	if _, err := io.Copy(part, f2); err != nil {
		t.Fatalf("erro copiando PNG para multipart: %v", err)
	}
}
