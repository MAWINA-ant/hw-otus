package storage

import (
	"errors"
	"time"
)

var (
	ErrDateBusy    = errors.New("date or time is busy")
	ErrIdNotFound  = errors.New("event with id not found")
	ErrAddEvent    = errors.New("couldn't add new event")
	ErrBadInterval = errors.New("bad interval")
)

type Event struct {
	ID          string
	Title       string
	DateTime    time.Time
	Duration    time.Duration
	Description string
	UserID      string
	NotifyTime  time.Time
}

