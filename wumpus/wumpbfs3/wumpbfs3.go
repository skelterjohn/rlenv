package main

import (
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"go-glue.googlecode.com/hg/rlglue"
	"github.com/skelterjohn/rlalg/bfs3"
	"github.com/skelterjohn/rlenv/wumpus"
)

func main() {
	defer nicetrace.Print()
	cfg := bfs3.ConfigDefault()
	argcfg.LoadArgs(&cfg)
	cfg.Gamma = 1
	agent := wumpus.NewBFS3Agent(cfg)
	rlglue.LoadAgent(agent)
}
