package watcher

import "testing"

func TestMetrics_InitialZero(t *testing.T) {
	m := NewMetrics()
	snap := m.Snapshot()
	if snap.Checks != 0 || snap.Missed != 0 || snap.Alerted != 0 {
		t.Errorf("expected all zeros, got %+v", snap)
	}
}

func TestMetrics_RecordCheck(t *testing.T) {
	m := NewMetrics()
	m.RecordCheck()
	m.RecordCheck()
	snap := m.Snapshot()
	if snap.Checks != 2 {
		t.Errorf("expected Checks=2, got %d", snap.Checks)
	}
}

func TestMetrics_RecordMissed(t *testing.T) {
	m := NewMetrics()
	m.RecordMissed()
	snap := m.Snapshot()
	if snap.Missed != 1 {
		t.Errorf("expected Missed=1, got %d", snap.Missed)
	}
}

func TestMetrics_RecordAlerted(t *testing.T) {
	m := NewMetrics()
	m.RecordAlerted()
	m.RecordAlerted()
	m.RecordAlerted()
	snap := m.Snapshot()
	if snap.Alerted != 3 {
		t.Errorf("expected Alerted=3, got %d", snap.Alerted)
	}
}

func TestMetrics_SnapshotIndependence(t *testing.T) {
	m := NewMetrics()
	snap1 := m.Snapshot()
	m.RecordCheck()
	snap2 := m.Snapshot()
	if snap1.Checks != 0 {
		t.Errorf("snap1 should be unchanged, got Checks=%d", snap1.Checks)
	}
	if snap2.Checks != 1 {
		t.Errorf("snap2 should reflect new check, got Checks=%d", snap2.Checks)
	}
}
