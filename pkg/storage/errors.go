package storage

import "errors"

var (
	// ErrKeyNotFound is returned when the provided key is not found in the store
	ErrKeyNotFound = errors.New(`key not found`)
	// ErrKeyAlreadyExists is returned when an operation would overwrite an exising key in the store
	ErrKeyAlreadyExists = errors.New(`key already exists`)
)
