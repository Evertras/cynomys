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

type UDPSender struct {
	mu sync.RWMutex

	broadcastAddr net.UDPAddr
	conn          *net.UDPConn
	sendInterval  time.Duration

	fromAddr string
	toAddr   string

	sendData []byte

	sink metrics.Sink
}

func NewUDPSender(addr net.UDPAddr, sendInterval time.Duration, sink metrics.Sink, sendData []byte) *UDPSender {
	return &UDPSender{
		broadcastAddr: addr,
		sendInterval:  sendInterval,
		sink:          sink,
		sendData:      sendData,
	}
}

func (s *UDPSender) Send(data []byte) error {
	if s.conn == nil {
		c, err := net.DialUDP("udp4", nil, &s.broadcastAddr)

		if err != nil {
			return fmt.Errorf("net.DialUDP: %w", err)
		}

		s.conn = c
	}

	reply := make([]byte, 16)
	sent := time.Now()

	_, err := s.conn.Write(data)

	if err != nil {
		return fmt.Errorf("s.conn.Write: %w", err)
	}

	_, _, err = s.conn.ReadFromUDP(reply)

	if err != nil {
		return fmt.Errorf("waiting for ping reply s.conn.ReadFromUDP: %w", err)
	}

	if reply[0] != constants.PingReply[0] {
		return fmt.Errorf("ping reply was not %q: %q", constants.PingReply, string(reply))
	}

	latency := time.Since(sent)

	err = s.sink.SendRTTLatencyMeasurement(s.conn.LocalAddr().String(), s.conn.RemoteAddr().String(), latency)

	if err != nil {
		log.Printf("Failed to send latency measurement: %v", err)
	}

	return nil
}

func (s *UDPSender) Run() error {
	s.mu.RLock()
	sendUDPTo := s.broadcastAddr.String()
	sendInterval := s.sendInterval
	s.mu.RUnlock()

	for {
		err := s.Send(s.sendData)
		if err != nil {
			log.Printf("Failed to send to %q: %v", sendUDPTo, err)
		}
		time.Sleep(sendInterval)
	}
}
