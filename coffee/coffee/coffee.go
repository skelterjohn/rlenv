package main

import (
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"github.com/skelterjohn/rlenv/coffee"
	"go-glue.googlecode.com/hg/rlglue"
)

func main() {
	defer nicetrace.Print()
	cfg := coffee.ConfigDefault()
	argcfg.LoadArgs(&cfg)
	env := coffee.NewEnv(cfg)
	if err := rlglue.LoadEnvironment(env); err != nil {
		println(err.String())
	}
}
