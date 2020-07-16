package storage

import (
	"net/url"
	"testing"
)

func TestMemoryStore(t *testing.T) {
	m := NewMemoryStore()
	assertLen := func(expected int) {
		if storeLen := len(m.urls); storeLen != expected {
			t.Errorf("unexpected entries in memory store: got %d, want %d", storeLen, expected)
		}
	}
	assertURLForKey := func(key, expected string) {
		u, found := m.ShortURL(key)
		if !found {
			t.Errorf("url %s expected in store but not found", expected)
		}
		if u == nil {
			t.Errorf("url %s expected in store but is nil", expected)
			return
		}

		if u.String() != expected {
			t.Errorf("unexpected short url in memory store: got %s, want %s", u.String(), expected)
		}
	}

	assertLen(0)
	u, found := m.ShortURL("")
	if found || u != nil {
		t.Errorf("unexpected short url in memory store: %s", u.String())
	}

	const (
		url1 = "http://url1.com"
		url2 = "http://url2.com"
	)

	m.AddURL("a", mustMkURL(url1))
	assertLen(1)
	assertURLForKey("a", url1)

	m.AddURL("b", mustMkURL(url1))
	assertLen(2)
	assertURLForKey("a", url1)
	assertURLForKey("b", url1)

	m.AddURL("b", mustMkURL(url2))
	assertLen(2)
	assertURLForKey("a", url1)
	assertURLForKey("b", url2)

	m.AddURL("c", mustMkURL(url2))
	assertLen(3)
	assertURLForKey("a", url1)
	assertURLForKey("b", url2)
	assertURLForKey("c", url2)

	m.DeleteURL(mustMkURL(url2))
	assertLen(1)
	assertURLForKey("a", url1)

	m.DeleteURLByKey("a")
	assertLen(0)
}

func mustMkURL(str string) url.URL {
	u, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	return *u
}
