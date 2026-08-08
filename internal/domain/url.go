package domain

import "time"

type URL struct {
	ID        uint64
	UserID    uint64
	LongURL   string
	ExpiredAt time.Time
}
