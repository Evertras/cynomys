package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cucumber/godog"
)

const configFileLocation = "cyn-bdd-config.yaml"

func (t *testContext) aConfigurationFileThatContains(contents *godog.DocString) error {
	file, err := os.Create(configFileLocation)

	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", configFileLocation, err)
	}

	_, err = io.WriteString(file, contents.Content)

	if err != nil {
		return fmt.Errorf("failed to write contents: %w", err)
	}

	return nil
}
