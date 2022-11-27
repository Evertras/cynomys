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

func (t *testContext) startInBackground(command string, args ...string) error {
	cmd, err := captured.StartInBackground(t.execCtx, command, args...)

	if err != nil {
		return fmt.Errorf("t.startCmd: %w", err)
	}

	t.cmds = append(t.cmds, cmd)

	return nil
}

func (t *testContext) startCynInBackground(args ...string) error {
	return t.startInBackground("../bin/cyn", args...)
}

func (t *testContext) waitSeconds(seconds int) error {
	time.Sleep(time.Second * time.Duration(seconds))

	return nil
}

func (t *testContext) thereIsNoOutput() error {
	for _, cmd := range t.cmds {
		stdout := cmd.Stdout()

		if len(stdout) > 0 {
			return fmt.Errorf("stdout output:\n%s", stdout)
		}

		stderr := cmd.Stderr()

		if len(stderr) > 0 {
			return fmt.Errorf("stderr output:\n%s", stderr)
		}
	}

	return nil
}

func (t *testContext) someStdoutContains(output string) error {
	for _, cmd := range t.cmds {
		stdout := cmd.Stdout()

		if strings.Contains(stdout, output) {
			return nil
		}
	}

	// Show any stderr for easier debugging
	for _, cmd := range t.cmds {
		stderr := cmd.Stderr()

		if len(stderr) > 0 {
			fmt.Println("vv STDERR vv")
			fmt.Println(stderr)
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
