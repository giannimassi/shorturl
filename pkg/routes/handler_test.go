package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const redirectTo = "https://example.com/"

// mockProviderl is a simple mock providing the short-url for the redirect handler
type mockProvider struct {
	urls map[string]string
}

func newMockProvider(urls map[string]string) *mockProvider {
	return &mockProvider{
		urls: urls,
	}
}

func (s *mockProvider) ShortURL(key string) (*url.URL, bool) {
	str, found := s.urls[key]
	if !found {
		return nil, false
	}
	u, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	return u, found
}

func (s *mockProvider) AddURL(key string, u url.URL) {
	s.urls[key] = u.String()
}

// DeleteURLByKey allows to remove a key-url association for the specified key
func (s *mockProvider) DeleteURLByKey(key string) {
	delete(s.urls, key)
}

// DeleteURL allows to remove all key-url association for the specified url
func (s *mockProvider) DeleteURL(u url.URL) {
	for k, v := range s.urls {
		if v == u.String() {
			delete(s.urls, k)
		}
	}
}

func Test_redirectHandler(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		urlMap map[string]string

		expectedStatusCode int
		expectedRedirectTo string
	}{
		{
			name: "ok/a",
			url:  "http://shorturl.com/a",

			expectedStatusCode: 301,
			expectedRedirectTo: "https://example.org/a",
		},
		{
			name: "ok/b",
			url:  "http://shorturl.com/b",

			expectedStatusCode: 301,
			expectedRedirectTo: "https://example.org/b",
		},
		{
			name: "ok/c",
			url:  "http://shorturl.com/c",

			expectedStatusCode: 301,
			expectedRedirectTo: "https://example.org/c",
		},
		{
			name: "ko/d",
			url:  "http://shorturl.com/d",

			expectedStatusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultURLMap := map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/b",
				"c": "https://example.org/c",
			}

			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			provider := newMockProvider(defaultURLMap)
			redirectHandler(provider).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Fatalf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}

			// We don't check for error since we only want to test if the location header is set
			// and we're not interested in testing the implementation of http.Response.
			if location, _ := w.Result().Location(); tt.expectedRedirectTo != "" {
				// check that location header is set correctly
				if location.String() != tt.expectedRedirectTo {

				}
			} else if location != nil {
				// location header should not be set
				t.Errorf("unexpected location header (should not be set): got %v", location.String())
			}

			expectedBody := fmt.Sprintf("<a href=\"%s\">Moved Permanently</a>.\n\n", tt.expectedRedirectTo)
			if tt.expectedStatusCode == 404 {
				expectedBody = ""
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
		key              string
		malformedPayload bool

		expectedStatusCode int
		expectedRedirectTo string
	}{
		{
			name:               "ok/a",
			key:                "a",
			expectedStatusCode: 200,
			expectedRedirectTo: "https://example.org/a",
		},
		{
			name: "ok/b",
			key:  "b",

			expectedStatusCode: 200,
			expectedRedirectTo: "https://example.org/b",
		},
		{
			name: "ok/c",
			key:  "c",

			expectedStatusCode: 200,
			expectedRedirectTo: "https://example.org/c",
		},
		{
			name: "ko/not found",
			key:  "d",

			expectedStatusCode: 404,
		},
		{
			name:             "ko/malformed",
			malformedPayload: true,

			expectedStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultURLMap := map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/b",
				"c": "https://example.org/c",
			}
			buf := bytes.Buffer{}
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(&struct {
				Key string
			}{Key: tt.key}); err != nil {
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
			provider := newMockProvider(defaultURLMap)
			infoHandler(provider).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Fatalf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}
			if tt.expectedStatusCode != 200 {
				return
			}

			dec := json.NewDecoder(w.Body)
			bodyPayload := struct {
				Key string
				URL string
			}{}
			if err := dec.Decode(&bodyPayload); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if bodyPayload.URL != tt.expectedRedirectTo {
				t.Errorf("unexpected url in body: got %v want %v", bodyPayload.URL, tt.expectedRedirectTo)
			}
			if bodyPayload.Key != tt.key {
				t.Errorf("unexpected key in body: got %v want %v", bodyPayload.Key, tt.key)
			}
		})
	}
}

func Test_deleteURLHandler(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		malformedPayload bool

		expectedStatusCode int
		expURLMap          map[string]string
	}{
		{
			name:               "ok/a",
			url:                "https://example.org/a",
			expectedStatusCode: 200,
			expURLMap: map[string]string{
				"c": "https://example.org/c",
			},
		},
		{
			name:               "ok/a",
			url:                "https://example.org/c",
			expectedStatusCode: 200,
			expURLMap: map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/a",
			},
		},
		{
			name: "ko/nothing-to-delete",
			url:  "https://example.org/d",

			expectedStatusCode: 200,
			expURLMap: map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/a",
				"c": "https://example.org/c",
			},
		},
		{
			name: "ko/url/malformed",
			url:  string([]byte{0x7f}), // ascii control character is one case where url parsing fails

			expectedStatusCode: 400,
			expURLMap: map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/a",
				"c": "https://example.org/c",
			},
		},
		{
			name:             "ko/payload/malformed",
			malformedPayload: true,

			expectedStatusCode: 400,
			expURLMap: map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/a",
				"c": "https://example.org/c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultURLMap := map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/a",
				"c": "https://example.org/c",
			}

			buf := bytes.Buffer{}
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(&struct {
				URL string
			}{URL: tt.url}); err != nil {
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
			provider := newMockProvider(defaultURLMap)
			deleteURLHandler(provider).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Fatalf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}

			compareMaps(t, tt.expURLMap, provider.urls)
		})
	}
}

func Test_deleteURLByKeyHandler(t *testing.T) {
	tests := []struct {
		name             string
		key              string
		malformedPayload bool

		expectedStatusCode int
		expURLMap          map[string]string
	}{
		{
			name:               "ok/a",
			key:                "a",
			expectedStatusCode: 200,
			expURLMap: map[string]string{
				"b": "https://example.org/a",
				"c": "https://example.org/c",
			},
		},
		{
			name:               "ok/a",
			key:                "c",
			expectedStatusCode: 200,
			expURLMap: map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/a",
			},
		},
		{
			name: "ko/nothing-to-delete",
			key:  "d",

			expectedStatusCode: 200,
			expURLMap: map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/a",
				"c": "https://example.org/c",
			},
		},
		{
			name:             "ko/payload/malformed",
			malformedPayload: true,

			expectedStatusCode: 400,
			expURLMap: map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/a",
				"c": "https://example.org/c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultURLMap := map[string]string{
				"a": "https://example.org/a",
				"b": "https://example.org/a",
				"c": "https://example.org/c",
			}

			buf := bytes.Buffer{}
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(&struct {
				Key string
			}{Key: tt.key}); err != nil {
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
			provider := newMockProvider(defaultURLMap)
			deleteURLByKeyHandler(provider).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Fatalf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}

			compareMaps(t, tt.expURLMap, provider.urls)
		})
	}
}

func compareMaps(t *testing.T, want, got map[string]string) {
	if len(want) != len(got) {
		t.Errorf("maps length don't match: want %d, got %d", len(want), len(got))
	}
	for k, v := range want {
		vv, found := got[k]
		if !found {
			t.Errorf("maps values for key %s don't match: want %v, not found in %v", k, v, got)
			continue
		}
		if vv != v {
			t.Errorf("maps values for key %s don't match: want %v, got %v", k, v, vv)
			continue
		}
	}
}

func Test_addURL(t *testing.T) {
	tests := []struct {
		name             string
		key              string
		url              string
		malformedPayload bool

		expectedStatusCode int
	}{
		{
			name: "ok/a",
			key:  "a",
			url:  "https://example.org/a",

			expectedStatusCode: 200,
		},

		{
			name: "ok/b",
			key:  "b",
			url:  "https://example.org/b",

			expectedStatusCode: 200,
		},

		{
			name: "ok/c",
			key:  "c",
			url:  "https://example.org/c",

			expectedStatusCode: 200,
		},

		{
			name:             "ko/malformed-payload",
			key:              "d",
			url:              "https://example.org/c",
			malformedPayload: true,

			expectedStatusCode: 400,
		},

		{
			name: "malformed-url",
			key:  "d",
			url:  string([]byte{0x7f}), // ascii control character is one case where url parsing fails

			expectedStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			dec := json.NewEncoder(&buf)
			if err := dec.Encode(&struct {
				Key string
				URL string
			}{Key: tt.key, URL: tt.url}); err != nil {
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
			provider := newMockProvider(make(map[string]string))
			addURLHandler(provider).ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatusCode {
				t.Errorf("wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}

			if tt.expectedStatusCode != 200 {
				return
			}

			if len(provider.urls) != 1 {
				t.Errorf("unexpected urls added: %#v", provider.urls)
			}

			if v, found := provider.urls[tt.key]; !found || v != tt.url {
				t.Errorf("unexpected url added: got %v, want %v", v, tt.url)
			}
		})
	}

}

func TestOnlyIf(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		url     string
		usePOST bool

		expectedStatusCode int
	}{
		{
			name:   "get/ok",
			method: "GET",
			url:    "http://shorturl.com/fhsdbf",

			expectedStatusCode: 200,
		},

		{
			name:    "get/not-allowed",
			method:  "GET",
			url:     "http://shorturl.com/fhsdbf",
			usePOST: true,

			expectedStatusCode: 405,
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
			method := "GET"
			if tt.usePOST {
				method = "POST"
			}
			handler := onlyIf(method, noopHandler)
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
