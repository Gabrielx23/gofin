package main

import (
	"log"
	"os"

	"gofin/cmd/cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
