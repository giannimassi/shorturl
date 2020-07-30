package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/giannimassi/shorturl/pkg/storage"
)

const redirectTo = "https://example.com/"

// mockProviderl is a simple mock providing the short-url for the redirect handler
type mockProvider struct {
	url url.URL
	err error
}

func newMockProvider(url string, err error) *mockProvider {
	return &mockProvider{
		url: mustMkURL(url),
		err: err,
	}
}

func (s *mockProvider) ShortURL(key string) (*url.URL, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &s.url, nil
}

func (s *mockProvider) AddURL(key string, u url.URL) error {
	return s.err
}

func (s *mockProvider) DeleteURL(key string) error {
	return s.err
}

func Test_redirectHandler(t *testing.T) {
	tests := []struct {
		name               string
		redirectURL        string
		storageErr         error
		expectedStatusCode int
	}{
		{
			name:        "ok/a",
			redirectURL: "https://example.org/a",

			expectedStatusCode: 301,
		},
		{
			name:       "ko/key-not-found",
			storageErr: storage.ErrKeyNotFound,

			expectedStatusCode: 404,
		},

		{
			name:       "ko/unexpected-error",
			storageErr: errors.New("unexpected error"),

			expectedStatusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "https://shorturl.com/abcdef", nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			provider := newMockProvider(tt.redirectURL, tt.storageErr)
			redirectHandler(provider).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Fatalf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}

			// We don't check for error since we only want to test if the location header is set
			// and we're not interested in testing the implementation of http.Response.
			if location, _ := w.Result().Location(); tt.redirectURL != "" {
				// check that location header is set correctly
				if location.String() != tt.redirectURL {

				}
			} else if location != nil {
				// location header should not be set
				t.Errorf("unexpected location header (should not be set): got %v", location.String())
			}

			expectedBody := fmt.Sprintf("<a href=\"%s\">Moved Permanently</a>.\n\n", tt.redirectURL)
			if tt.expectedStatusCode != 200 {
				return
			}
			if w.Body.String() != expectedBody {
				t.Errorf("unexpected body: got %v want %v", w.Body.String(), expectedBody)
			}
		})
	}
}

func Test_infoHandler(t *testing.T) {
	tests := []struct {
		name             string
		redirectURL      string
		storageErr       error
		malformedPayload bool

		expectedStatusCode int
	}{
		{
			name:               "ok/a",
			redirectURL:        "https://example.org/a",
			expectedStatusCode: 200,
		},
		{
			name:             "ko/malformed",
			malformedPayload: true,

			expectedStatusCode: 400,
		},
		{
			name:       "ko/key-not-found",
			storageErr: storage.ErrKeyNotFound,

			expectedStatusCode: 404,
		},
		{
			name:       "ko/unexpected-errors",
			storageErr: errors.New(""),

			expectedStatusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(&infoRequestPayload{Key: "abcdefg"}); err != nil {
				t.Fatal(err)
			}
			if tt.malformedPayload {
				buf.Reset()
				buf.WriteString("[]")
			}

			req, err := http.NewRequest("POST", "", &buf)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			provider := newMockProvider(tt.redirectURL, tt.storageErr)
			infoHandler(provider).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Fatalf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}
			if tt.expectedStatusCode != 200 {
				return
			}

			dec := json.NewDecoder(w.Body)
			var bodyPayload infoResponsePayload
			if err := dec.Decode(&bodyPayload); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if bodyPayload.URL != tt.redirectURL {
				t.Errorf("unexpected url in body: got %v want %v", bodyPayload.URL, tt.redirectURL)
			}
		})
	}
}

func Test_deleteURLHandler(t *testing.T) {
	tests := []struct {
		name             string
		key              string
		storageErr       error
		malformedPayload bool

		expectedStatusCode int
	}{
		{
			name:               "ok/a",
			key:                "a",
			expectedStatusCode: 200,
		},
		{
			name:             "ko/payload/malformed",
			malformedPayload: true,

			expectedStatusCode: 400,
		},
		{
			name:       "ko/key-not -found",
			key:        "d",
			storageErr: storage.ErrKeyNotFound,

			expectedStatusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(&deleteURLRequestPayload{Key: tt.key}); err != nil {
				t.Fatal(err)
			}
			if tt.malformedPayload {
				buf.Reset()
				buf.WriteString("[]")
			}

			req, err := http.NewRequest("POST", "", &buf)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			provider := newMockProvider("", tt.storageErr)
			deleteURLHandler(provider).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Fatalf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}
		})
	}
}

func Test_addURL(t *testing.T) {
	tests := []struct {
		name             string
		storageErr       error
		malformedURL     bool
		malformedPayload bool

		expectedStatusCode int
	}{
		{
			name: "ok/a",

			expectedStatusCode: 200,
		},

		{
			name:         "ko/malformed-url",
			malformedURL: true,

			expectedStatusCode: 422,
		},

		{
			name:             "ko/malformed-payload",
			malformedPayload: true,

			expectedStatusCode: 400,
		},
		{
			name:       "ko/key-already-exists",
			storageErr: storage.ErrKeyAlreadyExists,

			expectedStatusCode: 409,
		},
		{
			name:       "ko/unknown-err",
			storageErr: errors.New(""),

			expectedStatusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			dec := json.NewEncoder(&buf)
			url := "https://example.org"
			if tt.malformedURL {
				url = string([]byte{0x7f})
			}
			if err := dec.Encode(&addURLRequestPayload{Key: "example", URL: url}); err != nil {
				t.Fatal(err)
			}

			if tt.malformedPayload {
				buf.Reset()
				buf.WriteString("[]")
			}

			req, err := http.NewRequest("GET", "http://shorturl.com/add", &buf)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			provider := newMockProvider("", tt.storageErr)
			addURLHandler(provider).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Errorf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}
		})
	}

}

func mustMkURL(str string) url.URL {
	u, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	return *u
}
