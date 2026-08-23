package internalhttp

import (
	"fmt"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type EventJSON struct {
	ID          string     `json:"id,omitempty"`
	Title       string     `json:"title"`
	DateTime    time.Time  `json:"date_time"`
	Duration    string     `json:"duration"`
	Description string     `json:"description,omitempty"`
	UserID      string     `json:"user_id"`
	NotifyTime  *time.Time `json:"notify_time,omitempty"`
}

type CreateEventResponse struct {
	ID string `json:"id"`
}

type EventsListResponse struct {
	Events []EventJSON `json:"events"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (e EventJSON) toStorageEvent() (*storage.Event, error) {
	duration, err := time.ParseDuration(e.Duration)
	if err != nil {
		return nil, fmt.Errorf("invalid duration %q: %w", e.Duration, err)
	}

	event := &storage.Event{
		Title:       e.Title,
		DateTime:    e.DateTime,
		Duration:    duration,
		Description: e.Description,
		UserID:      e.UserID,
	}

	if e.NotifyTime != nil {
		event.NotifyTime = *e.NotifyTime
	}

	if e.ID != "" {
		id, err := uuid.Parse(e.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid id %q: %w", e.ID, err)
		}
		event.ID = id
	}

	return event, nil
}

func toEventJSON(e *storage.Event) EventJSON {
	eventJson := EventJSON{
		ID:          e.ID.String(),
		Title:       e.Title,
		DateTime:    e.DateTime,
		Duration:    e.Duration.String(),
		Description: e.Description,
		UserID:      e.UserID,
	}

	if !e.NotifyTime.IsZero() {
		notifyTime := e.NotifyTime
		eventJson.NotifyTime = &notifyTime
	}

	return eventJson
}

func eventsToEventsJSON(events []*storage.Event) []EventJSON {
	result := make([]EventJSON, 0, len(events))
	for _, e := range events {
		result = append(result, toEventJSON(e))
	}
	return result
}
