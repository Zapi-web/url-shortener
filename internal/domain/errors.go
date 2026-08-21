package domain

import "errors"

var (
	ErrUrlNotFound  = errors.New("url not found")
	ErrInputisEmpty = errors.New("given input is empty")
	ErrInvalidInput = errors.New("invalid input")
	ErrQueueFull    = errors.New("queue is full")
	ErrQueueClosed  = errors.New("queue is closed")
)
