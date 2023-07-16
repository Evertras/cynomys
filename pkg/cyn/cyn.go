package cyn

import (
	"sync"

	"github.com/evertras/cynomys/pkg/httpserver"
	"github.com/evertras/cynomys/pkg/listener"
	"github.com/evertras/cynomys/pkg/sender"
)

type Cyn struct {
	mu sync.RWMutex

	tcpListeners []*listener.TCPListener
	udpListeners []*listener.UDPListener
	tcpSenders   []*sender.TCPSender
	udpSenders   []*sender.UDPSender

	httpServer *httpserver.Server
}

func New() *Cyn {
	return &Cyn{}
}

func (c *Cyn) AddTCPListener(tcpListener *listener.TCPListener) {
	c.mu.Lock()
	c.tcpListeners = append(c.tcpListeners, tcpListener)
	c.mu.Unlock()
}

func (c *Cyn) TCPListeners() []*listener.TCPListener {
	c.mu.RLock()
	defer c.mu.RUnlock()

	listeners := make([]*listener.TCPListener, len(c.tcpListeners))
	copy(listeners, c.tcpListeners)

	return listeners
}

func (c *Cyn) AddUDPListener(udpListener *listener.UDPListener) {
	c.mu.Lock()
	c.udpListeners = append(c.udpListeners, udpListener)
	c.mu.Unlock()
}

func (c *Cyn) UDPListeners() []*listener.UDPListener {
	c.mu.RLock()
	defer c.mu.RUnlock()

	listeners := make([]*listener.UDPListener, len(c.udpListeners))
	copy(listeners, c.udpListeners)

	return listeners
}

func (c *Cyn) AddTCPSender(tcpSender *sender.TCPSender) {
	c.mu.Lock()
	c.tcpSenders = append(c.tcpSenders, tcpSender)
	c.mu.Unlock()
}

func (c *Cyn) AddUDPSender(udpSender *sender.UDPSender) {
	c.mu.Lock()
	c.udpSenders = append(c.udpSenders, udpSender)
	c.mu.Unlock()
}

func (c *Cyn) AddHTTPServer(addr string) {
	server := httpserver.NewServer(addr, c)

	c.mu.Lock()
	c.httpServer = server
	c.mu.Unlock()
}
