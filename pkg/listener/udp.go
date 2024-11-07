package listener

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/evertras/cynomys/pkg/constants"
)

type UdpListener struct {
	mu   sync.RWMutex
	addr net.UDPAddr
	cfg  UdpConfig
}

type UdpConfig struct {
	Echo bool
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
