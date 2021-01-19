// stats is a thread-safe interface that provides
// statistical data about a crawler's operations
package stats

import (
	"log"
	"sync/atomic"
	"time"
)

type Stats struct {
	// total is the total number of URLs the crawler has processed
	total int64

	// pending is the number of ongoing requests
	pending int64

	// completed is the number of requests that were successfully
	// fetched and processed by the crawler
	completed int64

	// failures is the number of requests that failed somewhere in
	// the crawl pipeline
	failures int64

	// startTime is a record of when the crawler began crawling
	startTime time.Time

	// duration is the total time it took for the web crawler to crawl
	duration time.Duration
}

// NewStats creates a new stats structure
func NewStats() *Stats {
	return &Stats{
		total:     0,
		pending:   0,
		completed: 0,
		failures:  0,
	}
}

// RecordNewOperation updates the internal counters to indicate a
// new crawler operation
func (s *Stats) RecordNewOperation() {
	atomic.AddInt64(&s.pending, 1)
	atomic.AddInt64(&s.total, 1)
}

// RecordOperationFailure updates the internal counters to indicate a
// failed crawler operation
func (s *Stats) RecordOperationFailure() {
	atomic.AddInt64(&s.pending, -1)
	atomic.AddInt64(&s.failures, 1)
}

// RecordOperationCompletion updates the internal counters to indicate a
// completed crawler operation
func (s *Stats) RecordOperationCompletion() {
	atomic.AddInt64(&s.pending, -1)
	atomic.AddInt64(&s.completed, 1)
}

// RecordStartTime records the time the crawler began crawling
func (s *Stats) RecordStartTime() {
	// Only set the start time if it hasn't been set before
	if s.startTime.IsZero() {
		s.startTime = time.Now()
		return
	}
}

// RecordTotalDuration records the total time it took for the
// crawler to complete all operations
func (s *Stats) RecordTotalDuration() {
	// Only record the duration if the start time has been set
	if s.startTime.IsZero() {
		return
	}

	// Only record the duration if it hasn't been set before
	if s.duration == 0 {
		s.duration = time.Since(s.startTime)
	}
}

// Total returns a counter of the total number of URLs the crawler has processed
func (s *Stats) Total() int64 {
	return atomic.LoadInt64(&s.total)
}

// Pending returns a counter of the number of ongoing requests
func (s *Stats) Pending() int64 {
	return atomic.LoadInt64(&s.pending)
}

// Failures returns a counter of the number of requests that failed somewhere in
// the crawl pipeline
func (s *Stats) Failures() int64 {
	return atomic.LoadInt64(&s.failures)
}

// Completed returns a counter of the the number of requests that were successfully
// fetched and processed by the crawler
func (s *Stats) Completed() int64 {
	return atomic.LoadInt64(&s.completed)
}

// Duration returns the total time it took for the web crawler to crawl
func (s *Stats) Duration() time.Duration {
	return s.duration
}

// Print prints out the crawler's operation stats
func (s *Stats) Print() {
	log.Printf(
		"✨ Total: %v. Pending: %v. Completed: %v. Failed: %v",
		s.Total(),
		s.Pending(),
		s.Completed(),
		s.Failures(),
	)
}
