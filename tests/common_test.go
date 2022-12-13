package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/evertras/cynomys/tests/captured"
)

type testContext struct {
	execCtx context.Context
	cmds    []*captured.RunningCmd
}

func (t *testContext) iSendAUDPPacketContaining(data string, addressRaw string) error {
	addr, err := net.ResolveUDPAddr("udp", addressRaw)

	if err != nil {
		return fmt.Errorf("net.ResolveUDPAddr for %q: %w", addressRaw, err)
	}

	client, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		return fmt.Errorf("net.DialUDP for %q: %w", addressRaw, err)
	}

	defer client.Close()

	_, err = client.Write([]byte(data))

	if err != nil {
		return fmt.Errorf("client.Write for %q: %w", addressRaw, err)
	}

	// Just to make sure it got there...
	time.Sleep(time.Millisecond * 50)

	return nil
}
