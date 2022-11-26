package main

import (
	"fmt"
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
		return fmt.Errorf("not implemented yet")

	default:
		return fmt.Errorf("unexpected protocol %q", protocol)
	}
}
