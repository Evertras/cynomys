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
		for _, conn := range t.tcpConnections {
			if err := conn.Close(); err != nil {
				// TODO: ??
				panic(err)
			}
		}

		t.tcpConnections = nil

		cancelScenario()
	})

	sc.Step(`^a configuration file that contains:$`, t.aConfigurationFileThatContains)
	sc.Step(`^I run cyn (.*)$`, t.iRunCyn)
	sc.Step(`^cyn is run with no flags or config$`, t.cynIsRunWithoutFlagsOrConfig)
	sc.Step(`^cyn is run with an unknown flag$`, t.cynIsRunWithAnUnknownFlag)
	sc.Step(`^cyn is started with the config file$`, t.cynIsRunWithTheConfigFile)
	sc.Step(`^I wait (\d+) seconds?$`, t.waitSeconds)
	sc.Step(`^I wait a moment$`, t.waitAMoment)
	sc.Step(`^there is no output$`, t.thereIsNoOutput)
	sc.Step(`^I send a UDP packet containing "(.*)" to (.*)$`, t.iSendAUDPPacketContaining)
	sc.Step(`^I connect with TCP to (.*)$`, t.iConnectWithTCPTo)
	sc.Step(`^I send "(.*)" over my TCP connection$`, t.iSendOverMyTCPConnection)
	sc.Step(`^I disconnect my TCP connection$`, t.iDisconnectMyTCPConnection)
	sc.Step(`^some|the stdout contains "(.*)"$`, t.someStdoutContains)
	sc.Step(`^some|the stderr contains "(.*)"$`, t.someStderrContains)
	sc.Step(`^the stdout does not contain "(.*)"$`, t.noStdoutContains)
	sc.Step(`^I reset the output$`, t.iResetTheOutput)
	sc.Step(`^I stop process #(\d+)$`, t.iStopProcess)
}
