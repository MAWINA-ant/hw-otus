package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDateBusy    = errors.New("date or time is busy")
	ErrIDNotFound  = errors.New("event with id not found")
	ErrAddEvent    = errors.New("couldn't add new event")
	ErrBadInterval = errors.New("bad interval")
)

type Event struct {
	ID          uuid.UUID
	Title       string
	DateTime    time.Time
	Duration    time.Duration
	Description string
	UserID      string
	NotifyTime  time.Time
}
