package main

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	var scenarioCtx context.Context
	var cancelScenario context.CancelFunc

	t := testContext{}

	sc.BeforeScenario(func(sc *godog.Scenario) {
		scenarioCtx, cancelScenario = context.WithCancel(context.Background())

		t.execCtx = scenarioCtx
		t.cmds = nil
	})

	sc.AfterScenario(func(sc *godog.Scenario, err error) {
		cancelScenario()
	})

	sc.Step(`^cyn is listening for (UDP|TCP) on (.*)$`, t.cynIsListeningFor)
	sc.Step(`^I wait (\d+) seconds?$`, t.waitSeconds)
	sc.Step(`^there is no output$`, t.thereIsNoOutput)
	sc.Step(`^I send a UDP packet containing "(.*)" to (.*)$`, t.iSendAUDPPacketContaining)
	sc.Step(`^some|the stdout contains "(.*)"$`, t.someStdoutContains)
}
