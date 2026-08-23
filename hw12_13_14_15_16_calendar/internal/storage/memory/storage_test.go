package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	memoryStorage := New()

	testUUID := uuid.New()

	ctx := context.Background()

	const userID = "one"

	first := storage.Event{
		ID:          testUUID,
		Title:       "first",
		DateTime:    time.Date(2026, time.August, 3, 10, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "first test event",
		UserID:      userID,
		NotifyTime:  time.Date(2026, time.August, 3, 9, 0, 0, 0, time.UTC),
	}

	second := storage.Event{
		ID:          uuid.New(),
		Title:       "second",
		DateTime:    time.Date(2026, time.August, 3, 20, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "second test event",
		UserID:      userID,
		NotifyTime:  time.Date(2026, time.August, 3, 19, 0, 0, 0, time.UTC),
	}

	third := storage.Event{
		ID:          uuid.New(),
		Title:       "third",
		DateTime:    time.Date(2026, time.August, 4, 10, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "third test event",
		UserID:      userID,
		NotifyTime:  time.Date(2026, time.August, 4, 9, 0, 0, 0, time.UTC),
	}

	fourth := storage.Event{
		ID:          uuid.New(),
		Title:       "fourth",
		DateTime:    time.Date(2026, time.August, 12, 10, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "fourth test event",
		UserID:      userID,
		NotifyTime:  time.Date(2026, time.August, 12, 9, 0, 0, 0, time.UTC),
	}

	forEdit := storage.Event{
		ID:          testUUID,
		Title:       "",
		DateTime:    time.Date(2026, time.August, 24, 10, 0, 0, 0, time.UTC),
		Duration:    1 * time.Minute,
		Description: "test event",
		UserID:      userID,
		NotifyTime:  time.Date(2026, time.August, 24, 9, 0, 0, 0, time.UTC),
	}

	t.Run("empty storage", func(t *testing.T) {
		err := memoryStorage.EditEvent(ctx, testUUID, &forEdit)
		require.Equal(t, err, storage.ErrIDNotFound)
		err = memoryStorage.RemoveEvent(ctx, testUUID)
		require.Equal(t, err, storage.ErrIDNotFound)
	})

	t.Run("create event", func(t *testing.T) {
		firstUUID, err := memoryStorage.CreateEvent(ctx, &first)
		require.Nil(t, err)
		require.Equal(t, 1, len(memoryStorage.eventMap))
		require.Equal(t, testUUID, firstUUID)
		_, err = memoryStorage.CreateEvent(ctx, &second)
		require.Nil(t, err)
		require.Equal(t, 2, len(memoryStorage.eventMap))
		_, err = memoryStorage.CreateEvent(ctx, &third)
		require.Nil(t, err)
		require.Equal(t, 3, len(memoryStorage.eventMap))
		_, err = memoryStorage.CreateEvent(ctx, &fourth)
		require.Nil(t, err)
		require.Equal(t, 4, len(memoryStorage.eventMap))
	})

	t.Run("day event", func(t *testing.T) {
		result, err := memoryStorage.DayEvents(ctx, time.Date(2026, time.August, 3, 0, 0, 0, 0, time.UTC))
		require.Nil(t, err)
		require.Equal(t, 2, len(result))
	})

	t.Run("week event", func(t *testing.T) {
		result, err := memoryStorage.WeekEvents(ctx, time.Date(2026, time.August, 3, 0, 0, 0, 0, time.UTC))
		require.Nil(t, err)
		require.Equal(t, 3, len(result))
	})

	t.Run("month event", func(t *testing.T) {
		result, err := memoryStorage.MonthEvents(ctx, time.Date(2026, time.August, 3, 0, 0, 0, 0, time.UTC))
		require.Nil(t, err)
		require.Equal(t, 4, len(result))
	})

	t.Run("edit event", func(t *testing.T) {
		err := memoryStorage.EditEvent(ctx, testUUID, &forEdit)
		require.Nil(t, err)
		require.Equal(t, 4, len(memoryStorage.eventMap))
	})

	t.Run("remove event", func(t *testing.T) {
		err := memoryStorage.RemoveEvent(ctx, testUUID)
		require.Nil(t, err)
		require.Equal(t, 3, len(memoryStorage.eventMap))
	})

	t.Run("day event after remove", func(t *testing.T) {
		result, err := memoryStorage.DayEvents(ctx, time.Date(2026, time.August, 3, 0, 0, 0, 0, time.UTC))
		require.Nil(t, err)
		require.Equal(t, 1, len(result))
	})
}
