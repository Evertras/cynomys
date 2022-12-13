package main

import "time"

func (t *testContext) waitSeconds(seconds int) error {
	time.Sleep(time.Second * time.Duration(seconds))

	return nil
}

func (t *testContext) waitAMoment() error {
	// This is arbitrary, but useful for letting things settle... bump this up if
	// things get flakey, but they really shouldn't be flakey...
	time.Sleep(time.Millisecond * 200)

	return nil
}
