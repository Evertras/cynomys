package main

import (
	"context"
	"net"
	"os"

	"github.com/evertras/cynomys/tests/captured"
)

type testContext struct {
	execCtx context.Context
	cmds    []*captured.RunningCmd

	tcpConnections []*net.TCPConn

	envResets map[string]string
}

func newTestContext() *testContext {
	return &testContext{
		envResets: make(map[string]string),
	}
}

func (t *testContext) cleanup() error {
	for _, conn := range t.tcpConnections {
		if err := conn.Close(); err != nil {
			return err
		}
	}

	t.tcpConnections = nil

	for k, v := range t.envResets {
		if v == "" {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v)
		}
	}

	t.envResets = make(map[string]string)
	t.cmds = nil

	return nil
}

func (t *testContext) addEnvReset(key, value string) {
	if _, exists := t.envResets[key]; !exists {
		t.envResets[key] = value
	}
}
