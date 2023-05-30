package main

import (
	"os"
)

func (t *testContext) envVarIsSet(key, value string) error {
	t.addEnvReset(key, os.Getenv(key))

	os.Setenv(key, value)

	return nil
}
