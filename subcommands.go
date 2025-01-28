package main

import (
	"flag"
)

const (
	ColorRed   = "\u001b[31m"
	ColorGreen = "\u001b[32m"
	ColorReset = "\u001b[0m"
)

type Command struct {
	flagset *flag.FlagSet
}

type runner interface {
	Run() error
	Init([]string)
	Name() string
}
