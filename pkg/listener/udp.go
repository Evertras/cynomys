package listener

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type UDPListener struct {
	addr net.UDPAddr
}

func NewUDP(addr net.UDPAddr) *UDPListener {
	return &UDPListener{
		addr: addr,
	}
}

func (l *UDPListener) Listen() error {
	conn, err := net.ListenUDP("udp", &l.addr)
	if err != nil {
		return fmt.Errorf("net.ListenUDP: %w", err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		rlen, remote, err := conn.ReadFromUDP(buf)

		if err != nil {
			return fmt.Errorf("conn.ReadFromUDP: %w", err)
		}

		log.Printf("Read %d bytes from %v", rlen, remote)
		log.Printf("Received: %s", strings.ReplaceAll(string(buf), "\n", "\\n"))
	}
}
