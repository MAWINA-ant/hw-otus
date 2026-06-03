package app

import (
	"context"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Storage interface {
	CreateEvent(e *storage.Event) error
	EditEvent(ID string, e *storage.Event) error
	RemoveEvent(ID string) error
	DayEvents(d time.Time) ([]storage.Event, error)
	WeekEvents(d time.Time) ([]storage.Event, error)
	MonthEvents(d time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	event := &storage.Event{}
	event.ID = id
	event.Title = title
	return a.storage.CreateEvent(event)
}

func (a *App) EditEvent(ctx context.Context, id, title string) error {
	event := &storage.Event{}
	event.ID = id
	event.Title = title
	return a.storage.EditEvent(id, event)
}

func (a *App) RemoveEvent(ctx context.Context, id string) error {
	return a.storage.RemoveEvent(id)
}
