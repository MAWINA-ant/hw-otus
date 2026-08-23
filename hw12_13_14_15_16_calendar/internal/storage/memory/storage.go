package memorystorage

import (
	"context"
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
	mu          sync.RWMutex
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

func (s *Storage) CreateEvent(_ context.Context, e *storage.Event) (uuid.UUID, error) { //lint:unused
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := slices.IndexFunc(s.sortedUUIDs, func(pair kv) bool {
		return pair.Value.Equal(e.DateTime)
	})
	if idx != -1 {
		return uuid.Nil, storage.ErrDateBusy
	}
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	s.eventMap[e.ID] = e
	s.sortedUUIDs = append(s.sortedUUIDs, kv{e.ID, e.DateTime})
	sort.Slice(s.sortedUUIDs, func(i, j int) bool {
		return s.sortedUUIDs[i].Value.Before(s.sortedUUIDs[j].Value)
	})
	return e.ID, nil
}

func (s *Storage) EditEvent(_ context.Context, id uuid.UUID, e *storage.Event) error { //lint:unused
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.eventMap[id]; !ok {
		return storage.ErrIDNotFound
	}
	s.eventMap[id] = e
	return nil
}

func (s *Storage) RemoveEvent(_ context.Context, id uuid.UUID) error { //lint:unused
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.eventMap[id]; !ok {
		return storage.ErrIDNotFound
	}
	delete(s.eventMap, id)
	return nil
}

func (s *Storage) DayEvents(_ context.Context, d time.Time) ([]*storage.Event, error) { //lint:unused
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

func (s *Storage) WeekEvents(_ context.Context, d time.Time) ([]*storage.Event, error) { //lint:unused
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

func (s *Storage) MonthEvents(_ context.Context, d time.Time) ([]*storage.Event, error) { //lint:unused
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
