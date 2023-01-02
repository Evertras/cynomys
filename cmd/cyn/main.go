package main

import (
	"fmt"
	"log"
	"os"

	"github.com/evertras/cynomys/cmd/cyn/cmds"
)

func main() {
	log.SetOutput(os.Stdout)
	err := cmds.Execute()

	if err != nil {
		fmt.Println("Failed to run:", err)
		os.Exit(1)
	}
}
