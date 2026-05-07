package watcher

import (
	"testing"
	"time"
)

func makeRecord(missed, alerted []string) CheckRecord {
	return CheckRecord{
		At:      time.Now(),
		Missed:  missed,
		Alerted: alerted,
	}
}

func TestHistory_InitialEmpty(t *testing.T) {
	h := NewHistory(5)
	if h.Len() != 0 {
		t.Fatalf("expected 0 records, got %d", h.Len())
	}
}

func TestHistory_RecordAndSnapshot(t *testing.T) {
	h := NewHistory(10)
	r := makeRecord([]string{"backup"}, []string{"backup"})
	h.Record(r)

	snap := h.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 record, got %d", len(snap))
	}
	if snap[0].Missed[0] != "backup" {
		t.Errorf("unexpected missed job: %s", snap[0].Missed[0])
	}
}

func TestHistory_EvictsOldestWhenFull(t *testing.T) {
	h := NewHistory(3)
	for i := 0; i < 4; i++ {
		h.Record(makeRecord([]string{string(rune('a' + i))}, nil))
	}

	if h.Len() != 3 {
		t.Fatalf("expected 3 records after eviction, got %d", h.Len())
	}
	snap := h.Snapshot()
	if snap[0].Missed[0] != "b" {
		t.Errorf("expected oldest remaining to be 'b', got %s", snap[0].Missed[0])
	}
}

func TestHistory_SnapshotIsIndependent(t *testing.T) {
	h := NewHistory(5)
	h.Record(makeRecord([]string{"job1"}, nil))

	snap := h.Snapshot()
	snap[0].Missed[0] = "mutated"

	snap2 := h.Snapshot()
	if snap2[0].Missed[0] == "mutated" {
		t.Error("snapshot mutation affected internal state")
	}
}

func TestHistory_DefaultCapWhenZero(t *testing.T) {
	h := NewHistory(0)
	for i := 0; i < 105; i++ {
		h.Record(makeRecord(nil, nil))
	}
	if h.Len() != 100 {
		t.Errorf("expected cap of 100, got %d", h.Len())
	}
}
