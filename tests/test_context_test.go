package main

import (
	"context"
	"net"

	"github.com/evertras/cynomys/tests/captured"
)

type testContext struct {
	execCtx context.Context
	cmds    []*captured.RunningCmd

	tcpConnections []*net.TCPConn
}
