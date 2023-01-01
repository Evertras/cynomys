package main

import (
	"fmt"
	"net"
)

func (t *testContext) iConnectWithTCPTo(addressRaw string) error {
	addr, err := net.ResolveTCPAddr("tcp", addressRaw)

	if err != nil {
		return fmt.Errorf("net.ResolveTCPAddr for %q: %w", addressRaw, err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		return fmt.Errorf("net.DialTCP for %q: %w", addressRaw, err)
	}

	t.tcpConnections = append(t.tcpConnections, conn)

	return nil
}

func (t *testContext) iSendOverMyTCPConnection(data string) error {
	if len(t.tcpConnections) != 1 {
		return fmt.Errorf("can only disconnect if there is one connection made")
	}

	conn := t.tcpConnections[0]

	_, err := conn.Write([]byte(data))

	if err != nil {
		return fmt.Errorf("conn.Write: %w", err)
	}

	return nil
}

func (t *testContext) iDisconnectMyTCPConnection() error {
	if len(t.tcpConnections) != 1 {
		return fmt.Errorf("can only disconnect if there is one connection made")
	}

	conn := t.tcpConnections[0]
	t.tcpConnections = nil

	err := conn.Close()

	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}
