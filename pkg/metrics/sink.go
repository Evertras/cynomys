package metrics

import "time"

// Sink is an interface for sending latency measurements.
type Sink interface {
	SendRTTLatencyMeasurement(fromAddr, toAddr string, measurement time.Duration) error
}
