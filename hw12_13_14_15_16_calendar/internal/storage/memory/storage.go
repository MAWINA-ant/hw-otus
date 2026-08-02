package memorystorage

import (
	"errors"
	"sync"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

var ErrEventNotFound = errors.New("event not found")

type Storage struct {
	mu       sync.RWMutex //nolint:unused
	eventMap map[uuid.UUID]*storage.Event
}

func New() *Storage {
	return &Storage{
		mu:       sync.RWMutex{},
		eventMap: make(map[uuid.UUID]*storage.Event),
	}
}

func (s *Storage) CreateEvent(e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventMap[e.ID] = e
	return nil
}

func (s *Storage) EditEvent(id uuid.UUID, e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.eventMap[id]; !ok {
		return ErrEventNotFound
	}
	s.eventMap[e.ID] = e
	return nil
}

func (s *Storage) RemoveEvent(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.eventMap[id]; !ok {
		return ErrEventNotFound
	}
	delete(s.eventMap, id)
	return nil
}

func (s *Storage) DayEvents(d time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.Unlock()
	return []storage.Event{}, nil
}

func (s *Storage) WeekEvents(d time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.Unlock()
	return []storage.Event{}, nil
}

func (s *Storage) MonthEvents(d time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.Unlock()
	return []storage.Event{}, nil
}
