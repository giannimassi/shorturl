package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/giannimassi/shorturl/pkg/storage"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

//go:generate swag init -g ./handler.go -o ../../docs

// ShortURLProvider is the repository from which short url are fetched.
type ShortURLProvider interface {
	// ShortURL returns true if a short url is found for the provided key, false otherwise
	ShortURL(key string) (*url.URL, error)
	// AddURL allows to store a key-url association
	AddURL(key string, u url.URL) error
	// DeleteURLByKey allows to Delete a key-url association for the specified key
	DeleteURL(key string) error
}

// @title Shorturl API
// @version 0.1
// @description This is an url shortening service
// @termsOfService http://swagger.io/terms/

// @contact.name Gianni Massi
// @contact.url http://www.shorturl.com/support
// @contact.email support@shorturl.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// Start runs the server, setting up all required routes
func Start(s ShortURLProvider) error {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.NoRoute(gin.WrapF(redirectHandler(s)))

	api := r.Group("/api")
	api.GET("", gin.WrapF(infoHandler(s)))
	api.PUT("", gin.WrapF(addURLHandler(s)))
	api.DELETE("", gin.WrapF(deleteURLHandler(s)))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r.Run(":8080")
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

// infoRequestPayload godoc
type infoRequestPayload struct {
	Key string
}

// infoResponsePayload godoc
type infoResponsePayload struct {
	Key string
	URL string
}

// infoHandler implements a handler that returns information about the key-url association
// @Summary Return short URL info
// @Description Returns information about the short url association stored for the provided key
// @Accept  json
// @Produce  json
// @Param payload body infoRequestPayload true "Key for which the request is made"
// @Success 200 {object} infoResponsePayload
// @Failure 400 "Payload cannot be decoded"
// @Failure 404 "Key not found"
// @Failure 500 "The server has encountered an unknown error"
// @Router /api [get]
func infoHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var inputPayload infoRequestPayload
		if err := dec.Decode(&inputPayload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")

		shortURL, err := s.ShortURL(inputPayload.Key)
		if errors.Is(err, storage.ErrKeyNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		outputPayload := infoResponsePayload{
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

// addURLRequestPayload godoc
type addURLRequestPayload struct {
	Key string
	URL string
}

// addURLHandler returns an http.Handler that allows to add a key-url association
// @Summary Add short url
// @Description Adds a new key-url association
// @Accept json
// @Param payload body addURLRequestPayload true "Key-url association to add"
// @Success 200 "Key-url association added"
// @Failure 400 "Payload cannot be decoded"
// @Failure 422 "URL in the payload is malformed"
// @Failure 409 "A key-url association already exists for the provided key"
// @Failure 500 "The server has encountered an unknown error"
// @Router /api [put]
func addURLHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var payload addURLRequestPayload
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

// deleteURLRequestPayload godoc
type deleteURLRequestPayload struct {
	Key string
}

// deleteURLByKeyHandler returns an http.Handler that allows to delete a key-url association
// @Summary Delete short url
// @Description Deletes a key-url association
// @Accept json
// @Param payload body deleteURLRequestPayload true "Key-url association to delete"
// @Success 200 "Key-url association deleted"
// @Failure 400 "Payload cannot be decoded"
// @Failure 404 "Key-url association not found for key"
// @Failure 500 "The server has encountered an unknown error"
// @Router /api [delete]
func deleteURLHandler(s ShortURLProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var payload deleteURLRequestPayload
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
