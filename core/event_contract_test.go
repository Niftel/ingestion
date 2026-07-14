package core

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/praetordev/events"
)

type capturePublisher struct {
	event *events.JobEvent
}

func (p *capturePublisher) PublishJobEvent(event *events.JobEvent) error {
	copy := *event
	p.event = &copy
	return nil
}

func (p *capturePublisher) PublishLogChunk(*events.LogChunk) error { return nil }

func TestIngestEventsPreservesWireFields(t *testing.T) {
	authenticatedRunID := uuid.MustParse("11111111-2222-3333-4444-555555555555")
	bodyRunID := uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	host, task, play, stdout := "host-1", "Install package", "Configure web tier", "ok: [host-1]"
	timestamp := time.Date(2026, 7, 14, 12, 1, 0, 0, time.UTC)
	publisher := &capturePublisher{}
	service := NewIngestionService(nil, publisher, nil)

	err := service.IngestEvents(context.Background(), authenticatedRunID, []events.JobEvent{{
		ExecutionRunID: bodyRunID,
		UnifiedJobID:   4242,
		Seq:            17,
		EventType:      "TASK_OK",
		Timestamp:      timestamp,
		Host:           &host,
		TaskName:       &task,
		PlayName:       &play,
		StdoutSnippet:  &stdout,
		EventData:      json.RawMessage(`{"changed":false}`),
	}})
	if err != nil {
		t.Fatal(err)
	}

	got := publisher.event
	if got == nil {
		t.Fatal("event was not published")
	}
	if got.ExecutionRunID != authenticatedRunID {
		t.Fatalf("run id = %s, want authenticated URL id %s", got.ExecutionRunID, authenticatedRunID)
	}
	if got.UnifiedJobID != 4242 || got.Seq != 17 || got.EventType != "TASK_OK" || !got.Timestamp.Equal(timestamp) {
		t.Fatalf("lifecycle fields were not preserved: %+v", got)
	}
	if got.Host == nil || *got.Host != host || got.TaskName == nil || *got.TaskName != task || got.PlayName == nil || *got.PlayName != play {
		t.Fatalf("observability fields were not preserved: %+v", got)
	}
	if string(got.EventData) != `{"changed":false}` {
		t.Fatalf("event_data = %s", got.EventData)
	}
}
