package main

import (
	"fmt"
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"github.com/skelterjohn/rlenv/coffee"
	"go-glue.googlecode.com/hg/rlglue"
)

func main() {
	defer nicetrace.Print()
	cfg := coffee.ConfigDefault()
	argcfg.LoadArgs(&cfg)
	agent := coffee.NewOptAgent(cfg)
	err := rlglue.LoadAgent(agent)
	if err != nil {
		fmt.Println(err)
	}
}
