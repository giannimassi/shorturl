package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/giannimassi/shorturl/pkg/storage"
)

// ShortURLProvider is the repository from which short url are fetched.
type ShortURLProvider interface {
	// ShortURL returns true if a short url is found for the provided key, false otherwise
	ShortURL(key string) (*url.URL, error)
	// AddURL allows to store a key-url association
	AddURL(key string, u url.URL) error
	// DeleteURLByKey allows to Delete a key-url association for the specified key
	DeleteURL(key string) error
}

// Mux returns a new http handler with routes for adding urls and redirecting
func Mux(s ShortURLProvider) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", onlyIf("GET", log(redirectHandler(s))))
	mux.Handle("/api/add", onlyIf("POST", log(addURLHandler(s))))
	mux.Handle("/api/delete", onlyIf("POST", log(deleteURLHandler(s))))
	mux.Handle("/api/info", onlyIf("POST", log(infoHandler(s))))
	return mux
}

// redirectHandler implements a handler that redirects to the url associated with the provided code
// NOTE: only GET requests are supported and tested.
// Reference: https://tools.ietf.org/html/rfc7231#section-6.4.2
func redirectHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := keyFromRequestURLPath(r.URL.Path)
		shortURL, err := s.ShortURL(key)
		if errors.Is(err, storage.ErrKeyNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
		shortURL, err := s.ShortURL(inputPayload.Key)
		if errors.Is(err, storage.ErrKeyNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if err := s.AddURL(payload.Key, *u); errors.Is(err, storage.ErrKeyAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			// TODO: return descriptive payload
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

// deleteURLByKeyHandler returns an http.Handler that allows to delete a key-url association
func deleteURLHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		payload := struct {
			Key string
		}{}
		if err := dec.Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := s.DeleteURL(payload.Key); errors.Is(err, storage.ErrKeyNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

func log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s - %s\n", r.Method, r.URL.String())
		handler.ServeHTTP(w, r)
	})
}
