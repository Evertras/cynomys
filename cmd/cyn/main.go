package main

import (
	"fmt"
	"os"
)

func main() {
	err := rootCmd.Execute()

	if err != nil {
		fmt.Println("Failed to run:", err)
		os.Exit(1)
	}
}
