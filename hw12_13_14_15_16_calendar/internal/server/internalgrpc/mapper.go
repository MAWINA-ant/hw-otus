package internalgrpc

import (
	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/server/internalgrpc/eventpb"
	"github.com/MAWINA-ant/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoEvent(e *storage.Event) *eventpb.Event {
	if e == nil {
		return nil
	}

	pe := &eventpb.Event{
		Id:          e.ID.String(),
		Title:       e.Title,
		DateTime:    timestamppb.New(e.DateTime),
		Duration:    durationpb.New(e.Duration),
		Description: e.Description,
		UserId:      e.UserID,
	}

	if !e.NotifyTime.IsZero() {
		pe.NotifyTime = timestamppb.New(e.NotifyTime)
	}

	return pe
}

func toProtoEvents(events []*storage.Event) []*eventpb.Event {
	result := make([]*eventpb.Event, 0, len(events))
	for _, e := range events {
		result = append(result, toProtoEvent(e))
	}
	return result
}
func fromProtoEvent(pe *eventpb.Event) (*storage.Event, error) {
	e := &storage.Event{
		Title:       pe.GetTitle(),
		Description: pe.GetDescription(),
		UserID:      pe.GetUserId(),
	}

	if pe.GetDateTime() != nil {
		e.DateTime = pe.GetDateTime().AsTime()
	}

	if pe.GetDuration() != nil {
		e.Duration = pe.GetDuration().AsDuration()
	}

	if pe.GetNotifyTime() != nil {
		e.NotifyTime = pe.GetNotifyTime().AsTime()
	}

	if pe.GetId() != "" {
		id, err := uuid.Parse(pe.GetId())
		if err != nil {
			return nil, err
		}
		e.ID = id
	}

	return e, nil
}
