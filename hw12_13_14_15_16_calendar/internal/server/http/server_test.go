package internalhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

type stubLogger struct{}

func (stubLogger) Error(_ ...interface{}) {}
func (stubLogger) Warn(_ ...interface{})  {}
func (stubLogger) Info(_ ...interface{})  {}
func (stubLogger) Debug(_ ...interface{}) {}

func newTestServer() *Server {
	calendar := app.New(stubLogger{}, memorystorage.New())
	return NewServer(stubLogger{}, calendar)
}

func TestHelloHandler(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/hello?name=Otus", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "Hello, Otus!", rec.Body.String())
}

func TestCreateEditDeleteEvent(t *testing.T) {
	srv := newTestServer()
	handler := srv.Handler()

	eventJSON := EventJSON{
		Title:       "meeting",
		DateTime:    time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC),
		Duration:    "30m",
		Description: "sync up",
		UserID:      "user-1",
	}

	var created CreateEventResponse
	t.Run("create", func(t *testing.T) {
		body, err := json.Marshal(eventJSON)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))
		require.NotEmpty(t, created.ID)
	})

	t.Run("create with bad duration", func(t *testing.T) {
		bad := eventJSON
		bad.Duration = "not-a-duration"
		body, err := json.Marshal(bad)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("day listing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/day?date=2026-08-20", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var list EventsListResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &list))
		require.Len(t, list.Events, 1)
		require.Equal(t, created.ID, list.Events[0].ID)
		require.Equal(t, "meeting", list.Events[0].Title)
	})

	t.Run("update", func(t *testing.T) {
		updated := eventJSON
		updated.Title = "meeting (rescheduled)"
		updated.DateTime = time.Date(2026, time.August, 21, 10, 0, 0, 0, time.UTC)
		body, err := json.Marshal(updated)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, "/events/"+created.ID, bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code)

		req = httptest.NewRequest(http.MethodGet, "/events/day?date=2026-08-21", nil)
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		var list EventsListResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &list))
		require.Len(t, list.Events, 1)
		require.Equal(t, "meeting (rescheduled)", list.Events[0].Title)
	})

	t.Run("update unknown id", func(t *testing.T) {
		body, err := json.Marshal(eventJSON)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, "/events/00000000-0000-0000-0000-000000000000", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/"+created.ID, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code)

		req = httptest.NewRequest(http.MethodGet, "/events/day?date=2026-08-21", nil)
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		var list EventsListResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &list))
		require.Empty(t, list.Events)
	})

	t.Run("delete unknown id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/"+created.ID, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestListEventsMissingDateParam(t *testing.T) {
	srv := newTestServer()
	handler := srv.Handler()

	for _, path := range []string{"/events/day", "/events/week", "/events/month"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code, "path %s", path)
	}
}
