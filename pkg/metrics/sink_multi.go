package metrics

import (
	"errors"
	"time"
)

// MultiSink sends latency measurements to multiple sinks.
// Simplifies the process of sending to multiple sinks.
// If no sinks are given, it will safely do nothing.
type MultiSink struct {
	sinks []Sink
}

// NewMultiSink creates a new MultiSink ready to send to the given sinks.
// If no sinks are given, it will safely do nothing.
func NewMultiSink(sinks ...Sink) *MultiSink {
	return &MultiSink{
		sinks: sinks,
	}
}

// SendLatencyMeasurement sends a latency measurement to all sinks.
// Any errors are returned as a single joined error.
func (s *MultiSink) SendLatencyMeasurement(fromAddr, toAddr string, measurement time.Duration) error {
	errs := make([]error, 0, len(s.sinks))

	for _, sink := range s.sinks {
		if err := sink.SendLatencyMeasurement(fromAddr, toAddr, measurement); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
