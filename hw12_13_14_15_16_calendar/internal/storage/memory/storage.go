package memorystorage

import (
	"sync"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu       sync.RWMutex //nolint:unused
	eventMap map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:       sync.RWMutex{},
		eventMap: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(e *storage.Event) error {

	return nil
}

func (s *Storage) EditEvent(id string, e *storage.Event) error {
	return nil
}

func (s *Storage) RemoveEvent(id string) error {
	return nil
}

func (s *Storage) DayEvents(d time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

func (s *Storage) WeekEvents(d time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

func (s *Storage) MonthEvents(d time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}
