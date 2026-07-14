package handler

import (
	"os"
	"strings"
	"testing"
)

func TestDecodeJobEventBatchV1(t *testing.T) {
	f, err := os.Open("testdata/job_event_batch.v1.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	events, err := decodeJobEvents(f)
	if err != nil {
		t.Fatalf("decode canonical v1 batch: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("decoded %d events, want 1", len(events))
	}
	event := events[0]
	if event.ExecutionRunID.String() != "11111111-2222-3333-4444-555555555555" || event.UnifiedJobID != 4242 {
		t.Fatalf("event identity was not preserved: %+v", event)
	}
	if event.Seq != 17 || event.EventType != "TASK_OK" || event.Timestamp.IsZero() {
		t.Fatalf("event lifecycle fields were not preserved: %+v", event)
	}
	if event.Host == nil || *event.Host != "host-1" || event.TaskName == nil || *event.TaskName != "Install package" {
		t.Fatalf("event observability fields were not preserved: %+v", event)
	}
	if !strings.Contains(string(event.EventData), `"duration_ms": 42`) {
		t.Fatalf("event_data was not preserved: %s", event.EventData)
	}
}
