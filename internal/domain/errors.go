package domain

import "errors"

var (
	ErrUrlNotFound  = errors.New("url not found")
	ErrInputisEmpty = errors.New("given input is empty")
)
