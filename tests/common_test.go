package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/evertras/cynomys/tests/captured"
)

type testContext struct {
	execCtx context.Context
	cmds    []*captured.RunningCmd
}

func (t *testContext) waitSeconds(seconds int) error {
	time.Sleep(time.Second * time.Duration(seconds))

	return nil
}

func (t *testContext) thereIsNoOutput() error {
	for _, cmd := range t.cmds {
		stdout := cmd.Stdout()

		if len(stdout) > 0 {
			fmt.Println(stdout)
			return fmt.Errorf("stdout output length: %d", len(stdout))
		}

		stderr := cmd.Stderr()

		if len(stderr) > 0 {
			fmt.Println(stderr)
			return fmt.Errorf("stderr output length: %d", len(stderr))
		}
	}

	return nil
}

func (t *testContext) someStdoutContains(output string) error {
	for _, cmd := range t.cmds {
		stdout := cmd.Stdout()

		fmt.Println(stdout)

		if strings.Contains(stdout, output) {
			return nil
		}
	}

	return fmt.Errorf("failed to find %q in any output", output)
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