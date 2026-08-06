package memorystorage

import (
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type kv struct {
	Key   uuid.UUID
	Value time.Time
}

type Storage struct {
	mu          sync.RWMutex //nolint:unused
	eventMap    map[uuid.UUID]*storage.Event
	sortedUUIDs []kv
}

func New() *Storage {
	return &Storage{
		mu:          sync.RWMutex{},
		eventMap:    make(map[uuid.UUID]*storage.Event),
		sortedUUIDs: make([]kv, 0),
	}
}

func (s *Storage) CreateEvent(e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := slices.IndexFunc(s.sortedUUIDs, func(pair kv) bool {
		return pair.Value.Equal(e.DateTime)
	})
	if idx != -1 {
		return storage.ErrDateBusy
	}
	s.eventMap[e.ID] = e
	s.sortedUUIDs = append(s.sortedUUIDs, kv{e.ID, e.DateTime})
	sort.Slice(s.sortedUUIDs, func(i, j int) bool {
		return s.sortedUUIDs[i].Value.Before(s.sortedUUIDs[j].Value)
	})
	return nil
}

func (s *Storage) EditEvent(id uuid.UUID, e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.eventMap[id]; !ok {
		return storage.ErrIdNotFound
	}
	s.eventMap[e.ID] = e
	return nil
}

func (s *Storage) RemoveEvent(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.eventMap[id]; !ok {
		return storage.ErrIdNotFound
	}
	delete(s.eventMap, id)
	return nil
}

func (s *Storage) DayEvents(d time.Time) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	year, month, day := d.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, d.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	result := make([]*storage.Event, 0)
	for _, event := range s.eventMap {
		if (event.DateTime.After(startOfDay) || event.DateTime.Equal(startOfDay)) &&
			event.DateTime.Before(endOfDay) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) WeekEvents(d time.Time) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	year, month, day := d.Date()
	startOfWeek := time.Date(year, month, day, 0, 0, 0, 0, d.Location())
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)
	result := make([]*storage.Event, 0)
	for _, event := range s.eventMap {
		if (event.DateTime.After(startOfWeek) || event.DateTime.Equal(startOfWeek)) &&
			event.DateTime.Before(endOfWeek) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) MonthEvents(d time.Time) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	year, month, day := d.Date()
	startOfMonth := time.Date(year, month, day, 0, 0, 0, 0, d.Location())
	endOfMonth := startOfMonth.Add(30 * 24 * time.Hour)
	result := make([]*storage.Event, 0)
	for _, event := range s.eventMap {
		if (event.DateTime.After(startOfMonth) || event.DateTime.Equal(startOfMonth)) &&
			event.DateTime.Before(endOfMonth) {
			result = append(result, event)
		}
	}
	return result, nil
}
