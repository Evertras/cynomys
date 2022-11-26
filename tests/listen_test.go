package main

import (
	"fmt"

	"github.com/cucumber/godog"
)

func (t *testContext) cynIsListeningFor(protocol, addr string) error {
	switch protocol {
	case "UDP":
		cmd, err := t.startCmd("../bin/cyn")

		if err != nil {
			return fmt.Errorf("t.startCmd: %w", err)
		}

		t.cmds = append(t.cmds, cmd)

		return nil

	case "TCP":
		return godog.ErrPending

	default:
		return fmt.Errorf("unexpected protocol %q", protocol)
	}
}
