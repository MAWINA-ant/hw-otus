package internalhttp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var eventJson EventJSON
	if err := json.NewDecoder(r.Body).Decode(&eventJson); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	event, err := eventJson.toStorageEvent()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := s.app.CreateEvent(r.Context(), event)
	if err != nil {
		writeStorageError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateEventResponse{ID: id.String()})
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid event id: "+err.Error())
		return
	}

	var eventJson EventJSON
	if err := json.NewDecoder(r.Body).Decode(&eventJson); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	event, err := eventJson.toStorageEvent()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.app.EditEvent(r.Context(), id, event); err != nil {
		writeStorageError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid event id: "+err.Error())
		return
	}

	if err := s.app.RemoveEvent(r.Context(), id); err != nil {
		writeStorageError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listDayEventsHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDateParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := s.app.DayEvents(r.Context(), date)
	if err != nil {
		writeStorageError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, EventsListResponse{Events: eventsToEventsJSON(events)})
}

func (s *Server) listWeekEventsHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDateParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := s.app.WeekEvents(r.Context(), date)
	if err != nil {
		writeStorageError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, EventsListResponse{Events: eventsToEventsJSON(events)})
}

func (s *Server) listMonthEventsHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDateParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := s.app.MonthEvents(r.Context(), date)
	if err != nil {
		writeStorageError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, EventsListResponse{Events: eventsToEventsJSON(events)})
}

// parseDateParam reads the required "date" query parameter (YYYY-MM-DD).
func parseDateParam(r *http.Request) (time.Time, error) {
	raw := r.URL.Query().Get("date")
	if raw == "" {
		return time.Time{}, errors.New(`missing required query parameter "date" (expected format YYYY-MM-DD)`)
	}

	date, err := time.Parse(dateLayout, raw)
	if err != nil {
		return time.Time{}, errors.New(`invalid "date" query parameter: expected format YYYY-MM-DD`)
	}

	return date, nil
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// writeStorageError maps a domain/storage error to an HTTP status code.
func writeStorageError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrIDNotFound), errors.Is(err, sql.ErrNoRows):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, storage.ErrDateBusy):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, storage.ErrBadInterval), errors.Is(err, storage.ErrAddEvent):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
