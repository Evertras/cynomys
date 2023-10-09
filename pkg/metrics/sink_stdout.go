package metrics

import (
	"log"
	"time"
)

// SinkStdout sends latency measurements to stdout.  Useful for ad hoc checks.
type SinkStdout struct{}

// NewSinkStdout creates a new SinkStdout ready to print to stdout.
func NewSinkStdout() *SinkStdout {
	return &SinkStdout{}
}

// SendLatencyMeasurement sends a latency measurement to stdout.
func (s *SinkStdout) SendLatencyMeasurement(fromAddr, toAddr string, measurement time.Duration) error {
	log.Printf("%v -> %v latency: %v", fromAddr, toAddr, measurement)

	return nil
}
