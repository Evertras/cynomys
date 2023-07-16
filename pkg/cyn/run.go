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

		// Shadow to correctly capture the variable
		tcpListener := tcpListener
		eg.Go(tcpListener.Listen)
	}

	for _, udpListener := range c.udpListeners {
		listenOrSendCount++

		// Shadow to correctly capture the variable
		udpListener := udpListener
		eg.Go(udpListener.Listen)
	}

	for _, tcpSender := range c.tcpSenders {
		listenOrSendCount++
		eg.Go(tcpSender.Run)
	}

	for _, udpSender := range c.udpSenders {
		listenOrSendCount++
		eg.Go(udpSender.Run)
	}

	c.mu.RUnlock()

	if listenOrSendCount == 0 {
		return fmt.Errorf("no listeners or senders specified")
	}

	return eg.Wait()
}
