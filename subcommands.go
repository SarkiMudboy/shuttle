package main

import (
	"flag"
)

type Color string

const (
	ColorRed    Color = "\u001b[31m"
	ColorGreen  Color = "\u001b[32m"
	ColorBlue   Color = "\u001b[34m"
	ColorYellow Color = "\u001b[33m"
	ColorReset  Color = "\u001b[0m"
)

type Command struct {
	flagset *flag.FlagSet
}

type runner interface {
	Run() error
	Init([]string)
	Name() string
}
