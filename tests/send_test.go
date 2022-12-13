package main

import (
	"fmt"

	"github.com/cucumber/godog"
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
		return godog.ErrPending

	default:
		return fmt.Errorf("unexpected protocol: %q", protocol)
	}
}
