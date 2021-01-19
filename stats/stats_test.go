package stats

import (
	"sync"
	"testing"
)

func TestRecordNewOperation(t *testing.T) {
	s := NewStats()

	expected := 100

	wg := sync.WaitGroup{}
	wg.Add(expected)

	for i := 0; i < expected; i++ {
		go func() {
			s.RecordNewOperation()
			wg.Done()
		}()
	}

	wg.Wait()

	if s.Total() != int64(expected) {
		t.Fatalf("expected total operations to be %v, got %v", expected, s.Total())
	}

	if s.Pending() != int64(expected) {
		t.Fatalf("expected pending operations to be %v, got %v", expected, s.Pending())
	}
}

func TestRecordOperationFailure(t *testing.T) {
	s := NewStats()

	expected := 200

	// Set fake pending tasks
	s.pending = int64(expected)

	wg := sync.WaitGroup{}
	wg.Add(expected)

	for i := 0; i < expected; i++ {
		go func() {
			s.RecordOperationFailure()
			wg.Done()
		}()
	}

	wg.Wait()

	if s.Failures() != int64(expected) {
		t.Fatalf("expected failed operations to be %v, got %v", expected, s.Total())
	}

	if s.Pending() != 0 {
		t.Fatalf("expected 0 pending operations, got %v", s.Pending())
	}
}

func TestRecordOperationCompletion(t *testing.T) {
	s := NewStats()

	expected := 250

	// Set fake pending tasks
	s.pending = int64(expected)

	wg := sync.WaitGroup{}
	wg.Add(expected)

	for i := 0; i < expected; i++ {
		go func() {
			s.RecordOperationCompletion()
			wg.Done()
		}()
	}

	wg.Wait()

	if s.Completed() != int64(expected) {
		t.Fatalf("expected completed operations to be %v, got %v", expected, s.Completed())
	}

	if s.Pending() != 0 {
		t.Fatalf("expected 0 pending operations, got %v", s.Pending())
	}
}

func TestRecordStartTime(t *testing.T) {
	s := NewStats()

	s.RecordStartTime()

	if s.startTime.IsZero() {
		t.Fatalf("expected start time to be recorded, got %v", s.startTime)
	}
}

func TestRecordTotalDuration(t *testing.T) {
	s := NewStats()

	s.RecordStartTime()
	s.RecordTotalDuration()

	if s.duration == 0 {
		t.Fatalf("expected total duration to be set")
	}
}
