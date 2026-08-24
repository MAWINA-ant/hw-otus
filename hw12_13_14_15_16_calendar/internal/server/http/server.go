package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Server struct {
	httpServer *http.Server
	logger     Logger
	app        Application
	config     ServerConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type Logger interface {
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

type Application interface {
	CreateEvent(ctx context.Context, event *storage.Event) (uuid.UUID, error)
	EditEvent(ctx context.Context, id uuid.UUID, event *storage.Event) error
	RemoveEvent(ctx context.Context, id uuid.UUID) error
	DayEvents(ctx context.Context, d time.Time) ([]*storage.Event, error)
	WeekEvents(ctx context.Context, d time.Time) ([]*storage.Event, error)
	MonthEvents(ctx context.Context, d time.Time) ([]*storage.Event, error)
}

type LogEntry struct {
	IP         string
	Timestamp  time.Time
	Method     string
	Path       string
	Protocol   string
	StatusCode int
	Latency    time.Duration
	UserAgent  string
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
		config: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}
}

func NewServerWithConfig(logger Logger, app Application, config ServerConfig) *Server {
	return &Server{
		logger: logger,
		app:    app,
		config: config,
	}
}

// Handler builds the full HTTP routing tree wrapped with the logging middleware.
// It is exported for tests, which can exercise it directly via httptest without
// binding a real listener.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.helloHandler)
	mux.HandleFunc("/hello", s.helloHandler)

	mux.HandleFunc("POST /events", s.createEventHandler)
	mux.HandleFunc("PUT /events/{id}", s.updateEventHandler)
	mux.HandleFunc("DELETE /events/{id}", s.deleteEventHandler)
	mux.HandleFunc("GET /events/day", s.listDayEventsHandler)
	mux.HandleFunc("GET /events/week", s.listWeekEventsHandler)
	mux.HandleFunc("GET /events/month", s.listMonthEventsHandler)

	return loggingMiddleware(mux, s.logger)
}

func (s *Server) Start(ctx context.Context) error {
	handler := s.Handler()

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		s.logger.Info("Starting HTTP server on ", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server error:", err)
		}
	}()
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(shutdownCtx)
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	response := fmt.Sprintf("Hello, %s!", name)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response) // #nosec G705
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		host, _, err := net.SplitHostPort(xff)
		if err == nil {
			return host
		}
		return xff
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
