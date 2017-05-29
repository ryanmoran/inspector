package main

import (
	"log"
	"os"

	"github.com/ryanmoran/inspector/commands"
	"github.com/ryanmoran/inspector/flags"
	"github.com/ryanmoran/inspector/tiles"
)

func main() {
	stderr := log.New(os.Stderr, "", 0)

	var global struct {
		Help bool   `short:"h" long:"help" description:"prints this usage information" default:"false"`
		Path string `short:"p" long:"path" description:"path to the product file"`
	}

	args, err := flags.Parse(&global, os.Args[1:])
	if err != nil {
		stderr.Fatal(err)
	}

	globalFlagsUsage, err := flags.Usage(global)
	if err != nil {
		stderr.Fatal(err)
	}

	var command string
	if len(args) > 0 {
		command, args = args[0], args[1:]
	}

	if command == "" {
		command = "help"
	}

	productParser := tiles.NewParser(global.Path, os.Stdout)

	commandSet := commands.Set{}
	commandSet["help"] = commands.NewHelp(os.Stdout, globalFlagsUsage, commandSet)
	commandSet["deadweight"] = commands.NewDeadweight(productParser, os.Stdout)
	commandSet["pkg-dep"] = commands.NewPkgDep(productParser, os.Stdout)

	err = commandSet.Execute(command, args)
	if err != nil {
		stderr.Fatal(err)
	}
}
