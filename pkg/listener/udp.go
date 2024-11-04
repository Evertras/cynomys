package listener

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/evertras/cynomys/pkg/constants"
)

type UDPListener struct {
	mu   sync.RWMutex
	addr net.UDPAddr
}

func NewUDP(addr net.UDPAddr) *UDPListener {
	return &UDPListener{
		addr: addr,
	}
}

func (l *UDPListener) Addr() string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.addr.String()
}

func (l *UDPListener) Listen() error {
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
		log.Printf("Received: %s", strings.ReplaceAll(string(buf), "\n", "\\n"))
	}
}
