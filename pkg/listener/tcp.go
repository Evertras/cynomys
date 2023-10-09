package listener

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/evertras/cynomys/pkg/constants"
)

type TCPListener struct {
	mu   sync.RWMutex
	addr net.TCPAddr
}

func NewTCP(addr net.TCPAddr) *TCPListener {
	return &TCPListener{
		addr: addr,
	}
}

func (l *TCPListener) Addr() string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.addr.String()
}

func (l *TCPListener) Listen() error {
	l.mu.RLock()
	listener, err := net.ListenTCP("tcp", &l.addr)
	l.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("net.ListenTCP: %w", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			return fmt.Errorf("listener.Accept: %w", err)
		}

		log.Printf("TCP connected from %s", conn.RemoteAddr().String())

		go func() {
			buf := make([]byte, 1024)
			remote := conn.RemoteAddr().String()

			for {
				rlen, err := conn.Read(buf)

				if err == io.EOF {
					if rlen > 0 {
						log.Printf("Read %d bytes from %v", rlen, remote)
						log.Printf("Received: %s", strings.ReplaceAll(string(buf), "\n", "\\n"))
					}

					break
				} else if err != nil {
					log.Printf("failed to read: %s", err.Error())
					_ = conn.Close()
					break
				}

				_, err = conn.Write([]byte(constants.PingReply))

				if err != nil {
					log.Printf("failed to reply to ping: %s", err.Error())
					_ = conn.Close()
					break
				}

				log.Printf("Read %d bytes from %v", rlen, remote)
				log.Printf("Received: %s", strings.ReplaceAll(string(buf), "\n", "\\n"))
			}

			log.Printf("TCP disconnected from %s", remote)
		}()
	}
}
