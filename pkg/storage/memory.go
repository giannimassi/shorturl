package storage

import (
	"net/url"
	"sync"
)

// MemoryStore is a simple mock providing a memory-based storage of key-url associations
type MemoryStore struct {
	m    sync.RWMutex
	urls map[string]urlData
}

type urlData struct {
	url  url.URL
	hits int
}

// NewMemoryStore returns a new copy of MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		urls: make(map[string]urlData),
	}
}

// ShortURL returns a url and true if the key matches an entry, nil and false if otherwise
func (s *MemoryStore) ShortURL(key string) (*url.URL, error) {
	s.m.Lock()
	defer s.m.Unlock()
	u, found := s.urls[key]
	if !found {
		return nil, ErrKeyNotFound
	}
	u.hits++
	s.urls[key] = u
	return &u.url, nil
}

// AddURL adds a key-url association
func (s *MemoryStore) AddURL(key string, u url.URL) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, found := s.urls[key]; found {
		return ErrKeyAlreadyExists
	}

	s.urls[key] = urlData{url: u}
	return nil
}

// DeleteURL allows to remove a key-url association for the specified key
func (s *MemoryStore) DeleteURL(key string) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, found := s.urls[key]; !found {
		return ErrKeyNotFound
	}
	delete(s.urls, key)
	return nil
}

// ShortURLInfo returns true and the number of hits if a short url is found for the provided key, false otherwise
func (s *MemoryStore) ShortURLInfo(key string) (*url.URL, int, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	u, found := s.urls[key]
	if !found {
		return nil, 0, ErrKeyNotFound
	}
	return &u.url, u.hits, nil
}
