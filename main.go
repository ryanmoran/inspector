package main

import (
	"log"
	"os"

	"github.com/ryanmoran/inspector/commands"
	"github.com/ryanmoran/inspector/flags"
)

func main() {
	stdout := log.New(os.Stdout, "", 0)

	var global struct {
		Help bool `short:"h" long:"help" description:"prints this usage information" default:"false"`
	}

	args, err := flags.Parse(&global, os.Args[1:])
	if err != nil {
		stdout.Fatal(err)
	}

	globalFlagsUsage, err := flags.Usage(global)
	if err != nil {
		stdout.Fatal(err)
	}

	var command string
	if len(args) > 0 {
		command, args = args[0], args[1:]
	}

	if command == "" {
		command = "help"
	}

	commandSet := commands.Set{}
	commandSet["help"] = commands.NewHelp(os.Stdout, globalFlagsUsage, commandSet)

	err = commandSet.Execute(command, args)
	if err != nil {
		stdout.Fatal(err)
	}
}
