package main

import (
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"github.com/skelterjohn/rlenv/sysadmin"
	"go-glue.googlecode.com/hg/rlglue"
	"go-glue.googlecode.com/hg/rltools/discrete"
)

func main() {
	defer nicetrace.Print()
	cfg := sysadmin.ConfigDefault()
	argcfg.LoadArgs(&cfg)
	if false {
		env := sysadmin.New(cfg)
		if err := rlglue.LoadEnvironment(env); err != nil {
			println(err.String())
		}
	} else {
		mdp := sysadmin.NewSysMDP(cfg)
		env := discrete.NewMDPEnvNoReset(mdp, mdp.Task, (1<<uint(cfg.NumSystems))-1)
		if err := rlglue.LoadEnvironment(env); err != nil {
			println(err.String())
		}
	}
}
