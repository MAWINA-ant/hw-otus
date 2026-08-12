package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"

	_ "github.com/jackc/pgx/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	dsn := "postgresql://user:password@localhost:5432/calendar?sslmode=disable"
	var err error
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	err = s.db.PingContext(ctx)
	if err != nil {
		return err
	}
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id UUID PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		date_time TIMESTAMP NOT NULL,
		duration BIGINT NOT NULL,
		description TEXT,
		user_id VARCHAR(255) NOT NULL,
		notify_time TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_events_date_time ON events(date_time);
	CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
	CREATE INDEX IF NOT EXISTS idx_events_notify_time ON events(notify_time);
	`

	_, err = s.db.ExecContext(ctx, query)
	return err
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) CreateEvent(e *storage.Event) error {
	query := `
	INSERT INTO events (id, title, date_time, duration, description, user_id, notify_time)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	e.ID = uuid.New()

	_, err := s.db.Exec(query, e.ID, e.Title, e.DateTime, int64(e.Duration), e.Description, e.UserID, e.NotifyTime)
	return err
}

func (s *Storage) EditEvent(id uuid.UUID, e *storage.Event) error {
	query := `
	UPDATE events 
	SET title = $1, date_time = $2, duration = $3, description = $4, user_id = $5, notify_time = $6
	WHERE id = $7
	`

	result, err := s.db.Exec(query, e.Title, e.DateTime, int64(e.Duration), e.Description, e.UserID, e.NotifyTime, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Storage) RemoveEvent(id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Storage) DayEvents(d time.Time) ([]*storage.Event, error) {
	startOfDay := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return s.scanEvents(startOfDay, endOfDay)
}

// WeekEvents возвращает события за неделю, начиная с указанного дня
func (s *Storage) WeekEvents(d time.Time) ([]*storage.Event, error) {
	// Находим начало недели (понедельник)
	weekday := int(d.Weekday())
	if weekday == 0 { // Воскресенье
		weekday = 7
	}
	startOfWeek := d.AddDate(0, 0, -(weekday - 1))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, d.Location())
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)

	return s.scanEvents(startOfWeek, endOfWeek)
}

func (s *Storage) MonthEvents(d time.Time) ([]*storage.Event, error) {
	startOfMonth := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return s.scanEvents(startOfMonth, endOfMonth)
}

func (s *Storage) scanEvents(startTime, endTime time.Time) ([]*storage.Event, error) {
	var events []*storage.Event

	query := `
	SELECT id, title, date_time, duration, description, user_id, notify_time
	FROM events
	WHERE date_time >= $1 AND date_time < $2
	ORDER BY date_time ASC
	`
	rows, err := s.db.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := &storage.Event{}
		var duration int64
		err := rows.Scan(&e.ID, &e.Title, &e.DateTime, &duration, &e.Description, &e.UserID, &e.NotifyTime)
		if err != nil {
			return nil, err
		}
		e.Duration = time.Duration(duration)
		events = append(events, e)
	}

	return events, nil
}
