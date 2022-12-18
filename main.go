package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/hashibuto/artillery"
	"github.com/hashibuto/shellfire/internal/command"
)

//go:embed VERSION
var VERSION string

func main() {
	command.Version = VERSION

	// Strip off the command itself
	args := os.Args[1:]

	processor := artillery.NewProcessor()
	processor.RemoveBuiltins(false)

	commands := []*artillery.Command{
		command.EvalCommand,
		command.PatternCommand,
		command.PayloadCommand,
		command.StripeCommand,
		command.VersionCommand,
	}

	for _, c := range commands {
		err := processor.AddCommand(c)
		if err != nil {
			panic(err)
		}
	}

	err := processor.Process(args)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
