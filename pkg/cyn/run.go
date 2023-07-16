package cyn

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

func (c *Cyn) Run() error {
	eg := errgroup.Group{}

	listenOrSendCount := 0

	c.mu.RLock()

	if c.httpServer != nil {
		eg.Go(c.httpServer.ServeAndListen)
	}

	for _, tcpListener := range c.tcpListeners {
		listenOrSendCount++
		listen := tcpListener.Listen
		eg.Go(listen)
	}

	for _, udpListener := range c.udpListeners {
		listenOrSendCount++
		listen := udpListener.Listen
		eg.Go(listen)
	}

	for _, tcpSender := range c.tcpSenders {
		listenOrSendCount++
		run := tcpSender.Run
		eg.Go(run)
	}

	for _, udpSender := range c.udpSenders {
		listenOrSendCount++
		run := udpSender.Run
		eg.Go(run)
	}

	c.mu.RUnlock()

	if listenOrSendCount == 0 {
		return fmt.Errorf("no listeners or senders specified")
	}

	return eg.Wait()
}
