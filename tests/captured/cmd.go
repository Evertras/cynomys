package captured

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
)

type capturedOutput struct {
	mu   sync.Mutex
	data []byte
}

func newCapturedOutput() *capturedOutput {
	return &capturedOutput{
		data: make([]byte, 0, 1024),
	}
}

func (c *capturedOutput) Write(data []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = append(c.data, data...)

	return len(data), nil
}

func (c *capturedOutput) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make([]byte, 0, 1024)
}

type RunningCmd struct {
	cmd    *exec.Cmd
	stdout *capturedOutput
	stderr *capturedOutput
}

func StartInBackground(ctx context.Context, command string, args ...string) (*RunningCmd, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	r := &RunningCmd{
		cmd:    cmd,
		stdout: newCapturedOutput(),
		stderr: newCapturedOutput(),
	}

	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr

	err := cmd.Start()

	if err != nil {
		return nil, fmt.Errorf("cmd.Start: %w", err)
	}

	return r, nil
}

func (r *RunningCmd) Stop() error {
	if r == nil || r.cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	err := r.cmd.Process.Kill()

	if err != nil {
		return fmt.Errorf("process.Kill: %w", err)
	}

	return nil
}

func (r *RunningCmd) Stdout() string {
	r.stdout.mu.Lock()
	defer r.stdout.mu.Unlock()

	output := string(r.stdout.data)

	return output
}

func (r *RunningCmd) Stderr() string {
	r.stderr.mu.Lock()
	defer r.stderr.mu.Unlock()

	output := string(r.stderr.data)

	return output
}

func (r *RunningCmd) ResetOutput() {
	r.stdout.reset()
	r.stderr.reset()
}
