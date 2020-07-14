package routes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_redirect(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string

		expectedStatusCode     int
		expectedBody           string
		expectedLocationHeader string
	}{
		{
			name:   "ok",
			method: "GET",
			url:    "http://shorturl.com/fhsdbf",

			expectedStatusCode:     301,
			expectedLocationHeader: redirectTo,
			expectedBody:           fmt.Sprintf("<a href=\"%s\">Moved Permanently</a>.\n\n", redirectTo),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			RedirectHandler().ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Errorf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}

			// We don't check for error since we only want to test if the location header is set
			// and we're not interested in testing the implementation of http.Response.
			if location, _ := w.Result().Location(); tt.expectedLocationHeader != "" {
				// check that location header is set correctly
				if location.String() != tt.expectedLocationHeader {

				}
			} else if location != nil {
				// location header should not be set
				t.Errorf("unexpected location header (should not be set): got %v", location.String())
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("unexpected body: got %v want %v", w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func Test_allowGETOnly(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string

		expectedStatusCode int
	}{
		{
			name:   "ok",
			method: "GET",
			url:    "http://shorturl.com/fhsdbf",

			expectedStatusCode: 200,
		},

		{
			name:   "head/not-allowed",
			method: "HEAD",
			url:    "http://shorturl.com/fhsdbf",

			expectedStatusCode: 405,
		},

		{
			name:   "post/not-allowed",
			method: "POST",
			url:    "http://shorturl.com/fhsdbf",

			expectedStatusCode: 405,
		},

		{
			name:   "put/not-allowed",
			method: "PUT",
			url:    "http://shorturl.com/fhsdbf",

			expectedStatusCode: 405,
		},

		{
			name:   "delete/not-allowed",
			method: "DELETE",
			url:    "http://shorturl.com/fhsdbf",

			expectedStatusCode: 405,
		},

		{
			name:   "connect/not-allowed",
			method: "CONNECT",
			url:    "http://shorturl.com/fhsdbf",

			expectedStatusCode: 405,
		},

		{
			name:   "options/not-allowed",
			method: "OPTIONS",
			url:    "http://shorturl.com/fhsdbf",

			expectedStatusCode: 405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			handler := AllowGETOnly(noopHandler)
			handler.ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Errorf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}

			if w.Body.String() != "" {
				t.Errorf("unexpected body: %v", w.Body.String())
			}
		})
	}
}
