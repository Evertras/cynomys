package main

import (
	"fmt"
	"strings"
)

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

	return fmt.Errorf("failed to find %q in any stdout output", output)
}

func (t *testContext) someStderrContains(output string) error {
	for _, cmd := range t.cmds {
		stderr := cmd.Stderr()

		if strings.Contains(stderr, output) {
			return nil
		}
	}

	// Show any stderr for easier debugging
	for _, cmd := range t.cmds {
		stdout := cmd.Stdout()

		if len(stdout) > 0 {
			fmt.Println("vv STDOUT vv")
			fmt.Println(stdout)
		}
	}

	return fmt.Errorf("failed to find %q in any stderr output", output)
}
