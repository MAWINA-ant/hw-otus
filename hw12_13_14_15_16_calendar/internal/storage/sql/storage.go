package sqlstorage

import (
	"context"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct { // TODO
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) CreateEvent(e *storage.Event) error {
	return nil
}

func (s *Storage) EditEvent(id uuid.UUID, e *storage.Event) error {
	return nil
}

func (s *Storage) RemoveEvent(id uuid.UUID) error {
	return nil
}

func (s *Storage) DayEvents(d time.Time) ([]*storage.Event, error) {
	return nil, nil
}

func (s *Storage) WeekEvents(d time.Time) ([]*storage.Event, error) {
	return nil, nil
}

func (s *Storage) MonthEvents(d time.Time) ([]*storage.Event, error) {
	return nil, nil
}
