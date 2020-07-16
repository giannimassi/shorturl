package storage

import (
	"net/url"
	"sync"
)

// MemoryStore is a simple mock providing a memory-based storage of key-url associations
type MemoryStore struct {
	m    sync.RWMutex
	urls map[string]url.URL
}

// NewMemoryStore returns a new copy of MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		urls: make(map[string]url.URL),
	}
}

// ShortURL returns a url and true if the key matches an entry, nil and false if otherwise
func (s *MemoryStore) ShortURL(key string) (*url.URL, bool) {
	s.m.RLock()
	defer s.m.RUnlock()
	u, found := s.urls[key]
	if !found {
		return nil, false
	}
	return &u, true
}

// AddURL adds a key-url association
func (s *MemoryStore) AddURL(key string, u url.URL) {
	s.m.Lock()
	defer s.m.Unlock()
	s.urls[key] = u
}

// DeleteURLByKey allows to remove a key-url association for the specified key
func (s *MemoryStore) DeleteURLByKey(key string) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.urls, key)
}

// DeleteURL allows to remove all key-url association for the specified url
func (s *MemoryStore) DeleteURL(u url.URL) {
	s.m.Lock()
	defer s.m.Unlock()
	// TODO: improve performance cost (e.g. keep url-keys association map `map[string][]string`)
	for k, v := range s.urls {
		if v == u {
			delete(s.urls, k)
		}
	}
}
