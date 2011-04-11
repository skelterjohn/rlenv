package main

import (
	"fmt"
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"github.com/skelterjohn/rlenv/sysadmin"
	"go-glue.googlecode.com/hg/rlglue"
)

func main() {
	defer nicetrace.Print()
	cfg := sysadmin.ConfigDefault()
	argcfg.LoadArgs(&cfg)
	agent := sysadmin.NewOptAgent(cfg)
	err := rlglue.LoadAgent(agent)
	if err != nil {
		fmt.Println(err)
	}
}
