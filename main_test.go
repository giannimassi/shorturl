package main

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
		// NOTE: all http methods should be tested here
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			http.HandlerFunc(redirect).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
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
