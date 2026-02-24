package main

import (
	"os"

	"github.com/geminik12/krag/cmd/app"
)

func main() {
	command := app.NewKragCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
