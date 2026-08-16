package app

import (
	"context"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

type Storage interface {
	CreateEvent(ctx context.Context, e *storage.Event) error
	EditEvent(ctx context.Context, ID uuid.UUID, e *storage.Event) error
	RemoveEvent(ctx context.Context, ID uuid.UUID) error
	DayEvents(ctx context.Context, d time.Time) ([]*storage.Event, error)
	WeekEvents(ctx context.Context, d time.Time) ([]*storage.Event, error)
	MonthEvents(ctx context.Context, d time.Time) ([]*storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		storage: storage,
		logger:  logger,
	}
}

func (a *App) CreateEvent(ctx context.Context, id uuid.UUID, title string) error {
	event := &storage.Event{}
	event.ID = id
	event.Title = title
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) EditEvent(ctx context.Context, id uuid.UUID, title string) error {
	event := &storage.Event{}
	event.ID = id
	event.Title = title
	return a.storage.EditEvent(ctx, id, event)
}

func (a *App) RemoveEvent(ctx context.Context, id uuid.UUID) error {
	return a.storage.RemoveEvent(ctx, id)
}

func (a *App) DayEvents(ctx context.Context, d time.Time) ([]*storage.Event, error) {
	return a.storage.DayEvents(ctx, d)
}

func (a *App) WeekEvents(ctx context.Context, d time.Time) ([]*storage.Event, error) {
	return a.storage.WeekEvents(ctx, d)
}

func (a *App) MonthEvents(ctx context.Context, d time.Time) ([]*storage.Event, error) {
	return a.storage.MonthEvents(ctx, d)
}
