package memorystorage

import (
	"testing"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	memoryStorage := New()

	testUUID := uuid.New()

	first := storage.Event{
		ID:          testUUID,
		Title:       "first",
		DateTime:    time.Date(2026, time.August, 3, 10, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "first test event",
		UserID:      "one",
		NotifyTime:  time.Date(2026, time.August, 3, 9, 0, 0, 0, time.UTC),
	}

	second := storage.Event{
		ID:          uuid.New(),
		Title:       "second",
		DateTime:    time.Date(2026, time.August, 3, 20, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "second test event",
		UserID:      "one",
		NotifyTime:  time.Date(2026, time.August, 3, 19, 0, 0, 0, time.UTC),
	}

	third := storage.Event{
		ID:          uuid.New(),
		Title:       "third",
		DateTime:    time.Date(2026, time.August, 4, 10, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "third test event",
		UserID:      "one",
		NotifyTime:  time.Date(2026, time.August, 4, 9, 0, 0, 0, time.UTC),
	}

	fourth := storage.Event{
		ID:          uuid.New(),
		Title:       "fourth",
		DateTime:    time.Date(2026, time.August, 12, 10, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "fourth test event",
		UserID:      "one",
		NotifyTime:  time.Date(2026, time.August, 12, 9, 0, 0, 0, time.UTC),
	}

	forEdit := storage.Event{
		ID:          testUUID,
		Title:       "",
		DateTime:    time.Date(2026, time.August, 24, 10, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "test event",
		UserID:      "one",
		NotifyTime:  time.Date(2026, time.August, 24, 9, 0, 0, 0, time.UTC),
	}

	t.Run("empty storage", func(t *testing.T) {
		err := memoryStorage.EditEvent(testUUID, &forEdit)
		require.Equal(t, err, storage.ErrIdNotFound)
		err = memoryStorage.RemoveEvent(testUUID)
		require.Equal(t, err, storage.ErrIdNotFound)
	})

	t.Run("create event", func(t *testing.T) {
		err := memoryStorage.CreateEvent(&first)
		require.Nil(t, err)
		require.Equal(t, 1, len(memoryStorage.eventMap))
		err = memoryStorage.CreateEvent(&second)
		require.Nil(t, err)
		require.Equal(t, 2, len(memoryStorage.eventMap))
		err = memoryStorage.CreateEvent(&third)
		require.Nil(t, err)
		require.Equal(t, 3, len(memoryStorage.eventMap))
		err = memoryStorage.CreateEvent(&fourth)
		require.Nil(t, err)
		require.Equal(t, 4, len(memoryStorage.eventMap))
	})

	t.Run("day event", func(t *testing.T) {
		result, err := memoryStorage.DayEvents(time.Date(2026, time.August, 3, 0, 0, 0, 0, time.UTC))
		require.Nil(t, err)
		require.Equal(t, 2, len(result))
	})

	t.Run("week event", func(t *testing.T) {
		result, err := memoryStorage.WeekEvents(time.Date(2026, time.August, 3, 0, 0, 0, 0, time.UTC))
		require.Nil(t, err)
		require.Equal(t, 3, len(result))
	})

	t.Run("month event", func(t *testing.T) {
		result, err := memoryStorage.MonthEvents(time.Date(2026, time.August, 3, 0, 0, 0, 0, time.UTC))
		require.Nil(t, err)
		require.Equal(t, 4, len(result))
	})

	t.Run("edit event", func(t *testing.T) {
		err := memoryStorage.EditEvent(testUUID, &forEdit)
		require.Nil(t, err)
		require.Equal(t, 4, len(memoryStorage.eventMap))
	})

	t.Run("remove event", func(t *testing.T) {
		err := memoryStorage.RemoveEvent(testUUID)
		require.Nil(t, err)
		require.Equal(t, 3, len(memoryStorage.eventMap))
	})

	t.Run("day event after remove", func(t *testing.T) {
		result, err := memoryStorage.DayEvents(time.Date(2026, time.August, 3, 0, 0, 0, 0, time.UTC))
		require.Nil(t, err)
		require.Equal(t, 1, len(result))
	})
}
