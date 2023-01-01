package main

import (
	"fmt"
)

func (t *testContext) cynIsListeningFor(protocol, addr string) error {
	switch protocol {
	case "UDP":
		err := t.startCynInBackground("--listen-udp", addr)

		if err != nil {
			return fmt.Errorf("t.startCynInBackground: %w", err)
		}

		return nil

	case "TCP":
		err := t.startCynInBackground("--listen-tcp", addr)

		if err != nil {
			return fmt.Errorf("t.startCynInBackground: %w", err)
		}

		return nil

	default:
		return fmt.Errorf("unexpected protocol %q", protocol)
	}
}
