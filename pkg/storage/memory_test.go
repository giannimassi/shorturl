package storage

import (
	"errors"
	"net/url"
	"testing"
)

func TestMemoryStore(t *testing.T) {
	m := NewMemoryStore()
	assertLen := func(expected int) {
		t.Helper()
		if storeLen := len(m.urls); storeLen != expected {
			t.Errorf("unexpected entries in memory store: got %d, want %d", storeLen, expected)
		}
	}
	assertURLForKey := func(key, expected string) {
		t.Helper()
		u, err := m.ShortURL(key)
		if err != nil {
			t.Errorf("err while getting short url for %v (expected %s): %v", key, expected, err)
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
	u, err := m.ShortURL("")
	if !errors.Is(err, errKeyNotFound) || u != nil {
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
	assertURLForKey("b", url1)

	m.AddURL("c", mustMkURL(url2))
	assertLen(3)
	assertURLForKey("a", url1)
	assertURLForKey("b", url1)
	assertURLForKey("c", url2)

	if err := m.DeleteURL(mustMkURL(url2)); err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	assertLen(2)
	assertURLForKey("a", url1)
	assertURLForKey("b", url1)

	if err := m.DeleteURLByKey("a"); err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	assertLen(1)

	if err := m.DeleteURL(mustMkURL(url1)); err != nil {

	}

	assertLen(0)

	if err := m.DeleteURLByKey("a"); !errors.Is(err, errKeyNotFound) {
		t.Errorf("unexpected err: %v", err)
	}
}

func mustMkURL(str string) url.URL {
	u, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	return *u
}
