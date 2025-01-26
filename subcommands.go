package main

import (
	"flag"
)

type Command struct {
	flagset *flag.FlagSet
}

type runner interface {
	Run() error
	Init([]string)
	Name() string
}
