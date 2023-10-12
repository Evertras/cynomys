package sender

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/evertras/cynomys/pkg/constants"
	"github.com/evertras/cynomys/pkg/metrics"
)

type TCPSender struct {
	mu sync.RWMutex

	broadcastAddr net.TCPAddr
	conn          *net.TCPConn
	sendInterval  time.Duration

	sink metrics.Sink
}

func NewTCPSender(addr net.TCPAddr, sendInterval time.Duration, sink metrics.Sink) *TCPSender {
	return &TCPSender{
		broadcastAddr: addr,
		sendInterval:  sendInterval,
		sink:          sink,
	}
}

func (s *TCPSender) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

func (s *TCPSender) Send(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		c, err := net.DialTCP("tcp4", nil, &s.broadcastAddr)

		if err != nil {
			return fmt.Errorf("net.DialTCP: %w", err)
		}

		s.conn = c
	}

	sent := time.Now()

	_, err := s.conn.Write(data)

	if err != nil {
		_ = s.conn.Close()
		s.conn = nil
		return fmt.Errorf("s.conn.Write: %w", err)
	}

	reply := make([]byte, 16)
	_, err = s.conn.Read(reply)

	if err != nil {
		return fmt.Errorf("waiting for ping reply s.conn.Read: %w", err)
	}

	if reply[0] != constants.PingReply[0] {
		return fmt.Errorf("ping reply was not %q: %q", constants.PingReply, string(reply))
	}

	latency := time.Since(sent)

	if err := s.sink.SendLatencyMeasurement(s.conn.LocalAddr().String(), s.conn.RemoteAddr().String(), latency); err != nil {
		log.Printf("Failed to send latency measurement: %v", err)
	}

	return nil
}

func (s *TCPSender) Run() error {
	s.mu.RLock()
	sendTCPTo := s.broadcastAddr.String()
	sendInterval := s.sendInterval
	s.mu.RUnlock()

	for {
		err := s.Send([]byte("hi"))

		if err != nil {
			log.Printf("Failed to send to %q: %v", sendTCPTo, err)
		}

		time.Sleep(sendInterval)
	}
}
