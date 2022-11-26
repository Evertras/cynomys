package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type runningCmd struct {
	cmd    *exec.Cmd
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

type testContext struct {
	execCtx context.Context
	cmds    []*runningCmd
}

func (t *testContext) startCmd(command string, args ...string) (*runningCmd, error) {
	cmd := exec.CommandContext(t.execCtx, command, args...)

	r := &runningCmd{
		cmd:    cmd,
		stdout: new(bytes.Buffer),
		stderr: new(bytes.Buffer),
	}

	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr

	err := cmd.Start()

	if err != nil {
		return nil, fmt.Errorf("cmd.Start: %w", err)
	}

	return r, nil
}

func (t *testContext) waitSeconds(seconds int) error {
	time.Sleep(time.Second * time.Duration(seconds))

	return nil
}

func (t *testContext) thereIsNoOutput() error {
	for _, cmd := range t.cmds {
		stdout := cmd.stdout.String()

		if len(stdout) > 0 {
			fmt.Println(stdout)
			return fmt.Errorf("stdout output length: %d", len(stdout))
		}

		stderr := cmd.stderr.String()

		if len(stderr) > 0 {
			fmt.Println(stderr)
			return fmt.Errorf("stderr output length: %d", len(stderr))
		}
	}

	return nil
}
