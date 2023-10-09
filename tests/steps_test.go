package main

import (
	"context"
	"fmt"
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

	t := newTestContext()

	sc.BeforeScenario(func(sc *godog.Scenario) {
		scenarioCtx, cancelScenario = context.WithCancel(context.Background())

		t.execCtx = scenarioCtx
	})

	sc.AfterScenario(func(sc *godog.Scenario, err error) {
		if err != nil {
			for i, cmd := range t.cmds {
				fmt.Printf("Process #%d\n", i)
				fmt.Println("vvvv STDOUT DUMP vvvv")
				fmt.Println(cmd.Stdout())
				fmt.Println("^^^^ STDOUT DUMP ^^^^")
				fmt.Println("")
				fmt.Println("vvvv STDERR DUMP vvvv")
				fmt.Println(cmd.Stderr())
				fmt.Println("^^^^ STDERR DUMP ^^^^")
			}
		}
		cleanupErr := t.cleanup()

		if cleanupErr != nil {
			// TODO: ??
			panic(cleanupErr)
		}

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
	sc.Step(`^the environment variable (.*) is set to "(.*)"$`, t.envVarIsSet)
	sc.Step(`the page at (.*) contains "(.*)"`, t.thePageAtContains)
}
