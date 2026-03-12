package domain

import "errors"

var (
	ErrUrlNotFound     = errors.New("url not found")
	ErrKeyAlreadyExist = errors.New("key already exists")
	ErrKeyisEmpty      = errors.New("given key is empty")
)
