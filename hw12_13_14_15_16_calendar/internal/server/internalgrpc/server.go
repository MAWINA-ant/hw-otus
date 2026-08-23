package internalgrpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/server/internalgrpc/eventpb"
	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Logger interface {
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

// Application is the subset of the calendar domain that the gRPC server needs.
type Application interface {
	CreateEvent(ctx context.Context, event *storage.Event) (uuid.UUID, error)
	EditEvent(ctx context.Context, id uuid.UUID, event *storage.Event) error
	RemoveEvent(ctx context.Context, id uuid.UUID) error
	DayEvents(ctx context.Context, d time.Time) ([]*storage.Event, error)
	WeekEvents(ctx context.Context, d time.Time) ([]*storage.Event, error)
	MonthEvents(ctx context.Context, d time.Time) ([]*storage.Event, error)
}

type ServerConfig struct {
	Host string
	Port int
}

type Server struct {
	eventpb.UnimplementedEventServiceServer

	grpcServer *grpc.Server
	logger     Logger
	app        Application
	config     ServerConfig
}

func NewServer(logger Logger, app Application, config ServerConfig) *Server {
	s := &Server{
		logger: logger,
		app:    app,
		config: config,
	}

	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(s.loggingInterceptor),
	)
	eventpb.RegisterEventServiceServer(s.grpcServer, s)

	return s
}

// Start starts the gRPC server and blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	go func() {
		s.logger.Info("Starting GRPC server on ", addr)
		if err := s.grpcServer.Serve(lsn); err != nil {
			s.logger.Error("GRPC server error: ", err)
		}
	}()

	<-ctx.Done()
	return nil
}

// Stop gracefully stops the gRPC server.
func (s *Server) Stop(_ context.Context) error {
	if s.grpcServer == nil {
		return nil
	}
	s.logger.Info("Shutting down GRPC server...")
	s.grpcServer.GracefulStop()
	return nil
}

func (s *Server) CreateEvent(c context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	event, err := fromProtoEvent(req.GetEvent())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.app.CreateEvent(c, event)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &eventpb.CreateEventResponse{Id: id.String()}, nil
}

func (s *Server) UpdateEvent(c context.Context, req *eventpb.UpdateEventRequest) (*eventpb.UpdateEventResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id: "+err.Error())
	}

	event, err := fromProtoEvent(req.GetEvent())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.app.EditEvent(c, id, event); err != nil {
		return nil, toGRPCError(err)
	}

	return &eventpb.UpdateEventResponse{}, nil
}

func (s *Server) DeleteEvent(c context.Context, req *eventpb.DeleteEventRequest) (*eventpb.DeleteEventResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id: "+err.Error())
	}

	if err := s.app.RemoveEvent(c, id); err != nil {
		return nil, toGRPCError(err)
	}

	return &eventpb.DeleteEventResponse{}, nil
}

func (s *Server) ListEventsForDay(
	ctx context.Context, req *eventpb.ListEventsRequest,
) (*eventpb.ListEventsResponse, error) {
	date, err := requireDate(req)
	if err != nil {
		return nil, err
	}

	events, err := s.app.DayEvents(ctx, date)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &eventpb.ListEventsResponse{Events: toProtoEvents(events)}, nil
}

func (s *Server) ListEventsForWeek(
	ctx context.Context, req *eventpb.ListEventsRequest,
) (*eventpb.ListEventsResponse, error) {
	date, err := requireDate(req)
	if err != nil {
		return nil, err
	}

	events, err := s.app.WeekEvents(ctx, date)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &eventpb.ListEventsResponse{Events: toProtoEvents(events)}, nil
}

func (s *Server) ListEventsForMonth(
	ctx context.Context, req *eventpb.ListEventsRequest,
) (*eventpb.ListEventsResponse, error) {
	date, err := requireDate(req)
	if err != nil {
		return nil, err
	}

	events, err := s.app.MonthEvents(ctx, date)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &eventpb.ListEventsResponse{Events: toProtoEvents(events)}, nil
}

func requireDate(req *eventpb.ListEventsRequest) (time.Time, error) {
	if req.GetDate() == nil {
		return time.Time{}, status.Error(codes.InvalidArgument, "date is required")
	}
	return req.GetDate().AsTime(), nil
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, storage.ErrIDNotFound), errors.Is(err, sql.ErrNoRows):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, storage.ErrDateBusy):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, storage.ErrBadInterval), errors.Is(err, storage.ErrAddEvent):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
