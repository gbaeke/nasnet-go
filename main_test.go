package main_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gbaeke/nasnet-go/handlers"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func init() {
	//instantiate route without certmagic SSL
	handlers.Routes()
}

func TestMainPage(t *testing.T) {
	t.Log("Given the need to test the main page.")
	{
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal("\tShould be able to create a request.", ballotX, err)
		}
		t.Log("\tShould be able to create a request.", checkMark)

		rw := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(rw, req)

		if rw.Code != 200 {
			t.Fatal("\tShould receive 200", ballotX, rw.Code)
		}
		t.Log("\tShould receive 200", checkMark)
	}

}

func TestInferPage(t *testing.T) {
	t.Log("Given the need to test the inference page.")
	{
		path := "./test/cat.jpg"
		file, err := os.Open(path)
		if err != nil {
			t.Fatal("\tShould be able to open a file.", ballotX, err)
		}
		defer file.Close()
		t.Log("\tShould be able to open a file.", checkMark)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("uploadfile", "cat.jpg")
		if err != nil {
			t.Fatal("\tShould be able to create form file.", ballotX, err)
		}
		t.Log("\tShould be able to create form file.", checkMark)

		_, err = io.Copy(part, file)
		err = writer.Close()

		req, err := http.NewRequest("POST", "/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		if err != nil {
			t.Fatal("\tShould be able to create request.", ballotX, err)
		}
		t.Log("\tShould be able to create request.", checkMark)

		rw := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(rw, req)

		result := rw.Body.String()
		if isCat := strings.Contains(result, "Egyptian_cat"); !isCat {
			t.Fatal("\tShould be able to infer Egyptian_cat.", ballotX, isCat)
		}
		t.Log("\tShould be able to infer Egyptian_cat.", checkMark)
	}

}
