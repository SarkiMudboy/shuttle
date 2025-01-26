package main

import (
	"fmt"
	"os"
)

func root(args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("Enter a subcommand")
	}

	command := args[0]
	request := NewRequest()

	subcommands := []runner{
		request,
	}

	for _, subcommand := range subcommands {
		if subcommand.Name() == command {
			subcommand.Init(args[1:])
			return subcommand.Run()
		}
	}

	return fmt.Errorf("Unknown command: %s", command)
}

func main() {

	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
