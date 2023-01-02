package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/evertras/cynomys/tests/captured"
)

func (t *testContext) cynIsRunWithoutFlagsOrConfig() error {
	return t.startCynInBackground()
}

func (t *testContext) cynIsRunWithAnUnknownFlag() error {
	return t.startCynInBackground("--thisdoesntexist")
}

func (t *testContext) xcynIsRunWithTheVersionCommand() error {
	return t.startCynInBackground("version")
}

func (t *testContext) cynIsRunWithTheConfigFile() error {
	return t.startCynInBackground("--config-file", configFileLocation)
}

func (t *testContext) iRunCyn(everything string) error {
	args := strings.Split(everything, " ")

	return t.startCynInBackground(args...)
}

func (t *testContext) startInBackground(command string, args ...string) error {
	cmd, err := captured.StartInBackground(t.execCtx, command, args...)

	if err != nil {
		return fmt.Errorf("t.startCmd: %w", err)
	}

	t.cmds = append(t.cmds, cmd)

	return nil
}

func (t *testContext) startCynInBackground(args ...string) error {
	err := t.startInBackground("../bin/cyn", args...)

	if err != nil {
		return err
	}

	// Give cyn a moment to actually start
	time.Sleep(time.Millisecond * 500)

	return nil
}
