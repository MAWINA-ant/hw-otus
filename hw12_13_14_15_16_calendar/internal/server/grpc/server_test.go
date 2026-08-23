package internalgrpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/app"
	eventpb "github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/server/grpc/pb"
	memorystorage "github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type stubLogger struct{}

func (stubLogger) Error(_ ...interface{}) {}
func (stubLogger) Warn(_ ...interface{})  {}
func (stubLogger) Info(_ ...interface{})  {}
func (stubLogger) Debug(_ ...interface{}) {}

func newTestClient(t *testing.T) eventpb.EventServiceClient {
	t.Helper()

	calendar := app.New(stubLogger{}, memorystorage.New())
	srv := NewServer(stubLogger{}, calendar, ServerConfig{})

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { _ = lis.Close() })

	go func() {
		_ = srv.grpcServer.Serve(lis)
	}()
	t.Cleanup(srv.grpcServer.Stop)

	dialer := func(context.Context, string) (net.Conn, error) { return lis.Dial() }

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	return eventpb.NewEventServiceClient(conn)
}

func TestCreateUpdateDeleteEvent(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t)

	dateTime := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC)

	createResp, err := client.CreateEvent(ctx, &eventpb.CreateEventRequest{
		Event: &eventpb.Event{
			Title:    "standup",
			DateTime: timestamppb.New(dateTime),
			Duration: durationpb.New(900),
			UserId:   "user-1",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, createResp.GetId())

	listResp, err := client.ListEventsForDay(ctx, &eventpb.ListEventsRequest{
		Date: timestamppb.New(dateTime),
	})
	require.NoError(t, err)
	require.Len(t, listResp.GetEvents(), 1)
	require.Equal(t, "standup", listResp.GetEvents()[0].GetTitle())

	_, err = client.UpdateEvent(ctx, &eventpb.UpdateEventRequest{
		Id: createResp.GetId(),
		Event: &eventpb.Event{
			Title:    "standup (moved)",
			DateTime: timestamppb.New(dateTime),
			Duration: durationpb.New(900),
			UserId:   "user-1",
		},
	})
	require.NoError(t, err)

	listResp, err = client.ListEventsForDay(ctx, &eventpb.ListEventsRequest{
		Date: timestamppb.New(dateTime),
	})
	require.NoError(t, err)
	require.Len(t, listResp.GetEvents(), 1)
	require.Equal(t, "standup (moved)", listResp.GetEvents()[0].GetTitle())

	_, err = client.DeleteEvent(ctx, &eventpb.DeleteEventRequest{Id: createResp.GetId()})
	require.NoError(t, err)

	listResp, err = client.ListEventsForDay(ctx, &eventpb.ListEventsRequest{
		Date: timestamppb.New(dateTime),
	})
	require.NoError(t, err)
	require.Empty(t, listResp.GetEvents())
}

func TestUpdateUnknownEventReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t)

	_, err := client.UpdateEvent(ctx, &eventpb.UpdateEventRequest{
		Id: "00000000-0000-0000-0000-000000000000",
		Event: &eventpb.Event{
			Title: "ghost",
		},
	})
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestListEventsRequiresDate(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t)

	_, err := client.ListEventsForDay(ctx, &eventpb.ListEventsRequest{})
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
