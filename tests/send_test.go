package main

import (
	"fmt"
)

func (t *testContext) cynIsSendingTo(protocol, addr string) error {
	switch protocol {
	case "UDP":
		err := t.startCynInBackground("--send-udp", addr)

		if err != nil {
			return fmt.Errorf("t.startCynInBackground: %w", err)
		}

		return nil

	case "TCP":
		err := t.startCynInBackground("--send-tcp", addr)

		if err != nil {
			return fmt.Errorf("t.startCynInBackground: %w", err)
		}

		return nil

	default:
		return fmt.Errorf("unexpected protocol: %q", protocol)
	}
}
