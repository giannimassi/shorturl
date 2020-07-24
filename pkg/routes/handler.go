package routes

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// ShortURLProvider is the repository from which short url are fetched.
type ShortURLProvider interface {
	// ShortURL returns true if a short url is found for the provided key, false otherwise
	ShortURL(key string) (*url.URL, bool)
	// AddURL allows to store a key-url association
	AddURL(key string, u url.URL)
	// DeleteURL allows to Delete all key-url association for the specified url
	DeleteURL(url url.URL) bool
	// DeleteURLByKey allows to Delete a key-url association for the specified key
	DeleteURLByKey(key string) bool
}

// Mux returns a new http handler with routes for adding urls and redirecting
func Mux(s ShortURLProvider) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", onlyIf("GET", redirectHandler(s)))
	mux.Handle("/add", onlyIf("POST", addURLHandler(s)))
	mux.Handle("/delete", onlyIf("POST", addURLHandler(s)))
	mux.Handle("/delete/bykey", onlyIf("POST", deleteURLByKeyHandler(s)))
	mux.Handle("/info", onlyIf("POST", infoHandler(s)))
	return mux
}

// redirectHandler implements a handler that redirects to the url associated with the provided code
// NOTE: only GET requests are supported and tested.
// Reference: https://tools.ietf.org/html/rfc7231#section-6.4.2
func redirectHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := keyFromRequestURLPath(r.URL.Path)
		shortURL, found := s.ShortURL(key)
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Redirect(w, r, shortURL.String(), http.StatusMovedPermanently)
	})
}

// infoHandler implements a handler that returns information about the key-url association
func infoHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		inputPayload := struct {
			Key string
		}{}
		if err := dec.Decode(&inputPayload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		shortURL, found := s.ShortURL(inputPayload.Key)
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		outputPayload := struct {
			Key string
			URL string
		}{
			Key: inputPayload.Key,
			URL: shortURL.String(),
		}
		if err := json.NewEncoder(w).Encode(&outputPayload); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func keyFromRequestURLPath(path string) string {
	return strings.TrimPrefix(path, "/")
}

// addURLHandler returns an http.Handler that allows to add a key-url association
func addURLHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		payload := struct {
			Key string
			URL string
		}{}
		if err := dec.Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u, err := url.Parse(payload.URL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.AddURL(payload.Key, *u)
	})
}

// deleteURLByKeyHandler returns an http.Handler that allows to delete a key-url association
func deleteURLHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		payload := struct {
			URL string
		}{}
		if err := dec.Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u, err := url.Parse(payload.URL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.DeleteURL(*u)
	})
}

// deleteURLByKeyHandler returns an http.Handler that allows to delete a key-url association
func deleteURLByKeyHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		payload := struct {
			Key string
		}{}
		if err := dec.Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.DeleteURLByKey(payload.Key)
	})
}

// Middlewares

// onlyIf is a middleware that responds with 405 to any method other than the one provided
func onlyIf(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
