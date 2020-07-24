package storage

import (
	"errors"
	"net/url"
	"sync"
)

// MemoryStore is a simple mock providing a memory-based storage of key-url associations
type MemoryStore struct {
	m    sync.RWMutex
	urls map[string]url.URL
}

var (
	errKeyNotFound      = errors.New(`key not found`)
	errKeyAlreadyExists = errors.New(`key already exists`)
	errURLNotFound      = errors.New(`url not found`)
)

// NewMemoryStore returns a new copy of MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		urls: make(map[string]url.URL),
	}
}

// ShortURL returns a url and true if the key matches an entry, nil and false if otherwise
func (s *MemoryStore) ShortURL(key string) (*url.URL, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	u, found := s.urls[key]
	if !found {
		return nil, errKeyNotFound
	}
	return &u, nil
}

// AddURL adds a key-url association
func (s *MemoryStore) AddURL(key string, u url.URL) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, found := s.urls[key]; found {
		return errKeyAlreadyExists
	}

	s.urls[key] = u
	return nil
}

// DeleteURLByKey allows to remove a key-url association for the specified key
func (s *MemoryStore) DeleteURLByKey(key string) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, found := s.urls[key]; !found {
		return errKeyNotFound
	}
	delete(s.urls, key)
	return nil
}

// DeleteURL allows to remove all key-url association for the specified url
func (s *MemoryStore) DeleteURL(u url.URL) error {
	s.m.Lock()
	defer s.m.Unlock()
	// TODO: improve performance cost (e.g. keep url-keys association map `map[string][]string`)
	var found bool
	for k, v := range s.urls {
		if v == u {
			found = true
			delete(s.urls, k)
		}
	}

	if !found {
		return errURLNotFound
	}

	return nil
}
