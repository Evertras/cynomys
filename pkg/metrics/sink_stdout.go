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

// SendRTTLatencyMeasurement sends a latency measurement to stdout.
func (s *SinkStdout) SendRTTLatencyMeasurement(fromAddr, toAddr string, measurement time.Duration) error {
	log.Printf("%v -> %v RTT latency: %v", fromAddr, toAddr, measurement)

	return nil
}
