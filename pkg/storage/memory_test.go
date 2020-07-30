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

	cmpURL := func(u *url.URL, err error, expected string, expectedErr error) {
		t.Helper()
		if err != expectedErr {
			t.Errorf("err %s expected in store but is %s", err, expectedErr)
		}

		if u == nil && expected != "" {
			t.Errorf("url %s expected in store but is nil", expected)
			return
		}

		if expected != "" && u.String() != expected {
			t.Errorf("unexpected short url in memory store: got %s, want %s", u.String(), expected)
		}
	}

	assertURLForKey := func(key, expected string) {
		t.Helper()
		u, err := m.ShortURL(key)
		cmpURL(u, err, expected, nil)
	}

	assertInfoForKey := func(key, expectedURL string, expectedHits int, expectedErr error) {
		t.Helper()
		u, hits, err := m.ShortURLInfo(key)
		cmpURL(u, err, expectedURL, expectedErr)
		if hits != expectedHits {
			t.Errorf("unexpected hits: got %d, want %d", hits, expectedHits)
		}
	}

	assertLen(0)
	u, err := m.ShortURL("")
	if !errors.Is(err, ErrKeyNotFound) || u != nil {
		t.Errorf("unexpected short url in memory store: %s", u.String())
	}

	const (
		url1 = "http://url1.com"
		url2 = "http://url2.com"
	)

	m.AddURL("a", mustMkURL(url1))
	assertLen(1)
	assertURLForKey("a", url1)
	assertInfoForKey("a", url1, 1, nil)

	m.AddURL("b", mustMkURL(url1))
	assertLen(2)
	assertURLForKey("a", url1)
	assertURLForKey("b", url1)
	assertInfoForKey("a", url1, 2, nil)
	assertInfoForKey("b", url1, 1, nil)

	m.AddURL("b", mustMkURL(url2))
	assertLen(2)
	assertURLForKey("a", url1)
	assertURLForKey("b", url1)
	assertInfoForKey("a", url1, 3, nil)
	assertInfoForKey("b", url1, 2, nil)

	m.AddURL("c", mustMkURL(url2))
	assertLen(3)
	assertURLForKey("a", url1)
	assertURLForKey("b", url1)
	assertURLForKey("c", url2)
	assertInfoForKey("a", url1, 4, nil)
	assertInfoForKey("b", url1, 3, nil)
	assertInfoForKey("c", url2, 1, nil)

	if err := m.DeleteURL("a"); err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	assertLen(2)
	assertInfoForKey("a", "", 0, ErrKeyNotFound)

	if err := m.DeleteURL("b"); err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	assertLen(1)

	if err := m.DeleteURL("c"); err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	assertLen(0)

	if err := m.DeleteURL("a"); !errors.Is(err, ErrKeyNotFound) {
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
