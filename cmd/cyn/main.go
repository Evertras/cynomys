package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	err := rootCmd.Execute()

	if err != nil {
		fmt.Println("Failed to run:", err)
		os.Exit(1)
	}
}
