package listener

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/evertras/cynomys/pkg/constants"
)

type UdpListener struct {
	mu   sync.RWMutex
	addr net.UDPAddr
	cfg  UdpConfig
}

type UdpConfig struct {
	Echo bool

	Burst BurstConfig
}

func NewUDP(addr net.UDPAddr, cfg UdpConfig) *UdpListener {
	return &UdpListener{
		addr: addr,
		cfg:  cfg,
	}
}

func (l *UdpListener) Addr() string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.addr.String()
}

func (l *UdpListener) Listen() error {
	l.mu.RLock()
	burst := l.cfg.Burst.Window != 0
	l.mu.RUnlock()

	if burst {
		return l.listenBurst()
	}

	return l.listenSingle()
}

func (l *UdpListener) listenSingle() error {
	l.mu.RLock()
	conn, err := net.ListenUDP("udp", &l.addr)
	l.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("net.ListenUDP: %w", err)
	}
	defer conn.Close()

	buf := make([]byte, maxPacketBytesReceived)

	for {
		rlen, remote, err := conn.ReadFromUDP(buf)

		if err != nil {
			return fmt.Errorf("conn.ReadFromUDP: %w", err)
		}

		// Write back before doing anything else to minimize latency,
		// every nanosecond counts!
		_, err = conn.WriteToUDP([]byte(constants.PingReply), remote)

		if err != nil {
			return fmt.Errorf("conn.WriteToUDP: %w", err)
		}

		log.Printf("Read %d bytes from %v", rlen, remote)

		if l.cfg.Echo {
			log.Printf("Received: %s", strings.ReplaceAll(string(buf), "\n", "\\n"))
		}
	}
}

func (l *UdpListener) listenBurst() error {
	l.mu.RLock()
	conn, err := net.ListenUDP("udp", &l.addr)
	cfg := l.cfg
	l.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("net.ListenUDP: %w", err)
	}
	defer conn.Close()

	data := make(chan []byte, 0)
	errs := make(chan error, 1)

	go func() {
		buf := make([]byte, maxPacketBytesReceived)

		for {
			rlen, _, err := conn.ReadFromUDP(buf)

			if err != nil {
				errs <- fmt.Errorf("conn.ReadFromUDP: %w", err)
				return
			}

			data <- buf[:rlen]
		}
	}()

	var aggregateData bytes.Buffer

	var timer *time.Timer

	for {
		select {
		case rcv := <-data:
			_, err := aggregateData.Write(rcv)
			if err != nil {
				return fmt.Errorf("failed to write to aggregate: %w", err)
			}

			timer = time.NewTimer(cfg.Burst.Window)

		case err := <-errs:
			return err
		}

	TIMEDLOOP:
		for {
			select {
			case rcv := <-data:
				_, err := aggregateData.Write(rcv)
				if err != nil {
					return fmt.Errorf("failed to write to aggregate: %w", err)
				}

			case err := <-errs:
				return err

			case <-timer.C:
				timer = nil

				log.Printf("Read %d bytes", aggregateData.Len())

				if l.cfg.Echo {
					log.Printf("Received: %s", strings.ReplaceAll(aggregateData.String(), "\n", "\\n"))
				}

				break TIMEDLOOP
			}
		}
	}
}
