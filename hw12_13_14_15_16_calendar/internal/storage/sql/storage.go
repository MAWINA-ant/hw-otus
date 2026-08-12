package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib" // PostgreSQL driver for database/sql
)

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

type StorageConfig struct {
	Host     string
	DataBase string
	User     string
	Password string
}

func (s *Storage) Connect(ctx context.Context) error {
	dsn := "postgresql://user:password@localhost:5432/calendar?sslmode=disable" //nolint:gosec
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

func (s *Storage) ConnectWithConfig(ctx context.Context, config StorageConfig) error {
	dsn := "postgresql://" + config.User + ":" + config.Password + "@" +
		config.Host + "/" + config.DataBase + "?sslmode=disable"
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

func (s *Storage) Close(_ context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, e *storage.Event) error {
	query := `
	INSERT INTO events (id, title, date_time, duration, description, user_id, notify_time)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	e.ID = uuid.New()

	_, err := s.db.ExecContext(ctx, query, e.ID, e.Title, e.DateTime, int64(e.Duration),
		e.Description, e.UserID, e.NotifyTime)
	return err
}

func (s *Storage) EditEvent(ctx context.Context, id uuid.UUID, e *storage.Event) error {
	query := `
	UPDATE events 
	SET title = $1, date_time = $2, duration = $3, description = $4, user_id = $5, notify_time = $6
	WHERE id = $7
	`

	result, err := s.db.ExecContext(ctx, query, e.Title, e.DateTime, int64(e.Duration),
		e.Description, e.UserID, e.NotifyTime, id)
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

func (s *Storage) RemoveEvent(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
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

func (s *Storage) DayEvents(ctx context.Context, d time.Time) ([]*storage.Event, error) {
	year, month, day := d.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, d.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return s.scanEvents(ctx, startOfDay, endOfDay)
}

func (s *Storage) WeekEvents(ctx context.Context, d time.Time) ([]*storage.Event, error) {
	year, month, day := d.Date()
	startOfWeek := time.Date(year, month, day, 0, 0, 0, 0, d.Location())
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)

	return s.scanEvents(ctx, startOfWeek, endOfWeek)
}

func (s *Storage) MonthEvents(ctx context.Context, d time.Time) ([]*storage.Event, error) {
	year, month, day := d.Date()
	startOfMonth := time.Date(year, month, day, 0, 0, 0, 0, d.Location())
	endOfMonth := startOfMonth.Add(30 * 24 * time.Hour)

	return s.scanEvents(ctx, startOfMonth, endOfMonth)
}

func (s *Storage) scanEvents(ctx context.Context, startTime, endTime time.Time) ([]*storage.Event, error) {
	var events []*storage.Event

	query := `
	SELECT id, title, date_time, duration, description, user_id, notify_time
	FROM events
	WHERE date_time >= $1 AND date_time < $2
	ORDER BY date_time ASC
	`
	rows, err := s.db.QueryContext(ctx, query, startTime, endTime)
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
