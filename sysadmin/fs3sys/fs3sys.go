package main

import (
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"go-glue.googlecode.com/hg/rlglue"
	"github.com/skelterjohn/rlbayes"
	"github.com/skelterjohn/rlalg/bfs3"
	"github.com/skelterjohn/rlenv/sysadmin"
)

func GetTruthFunc(cfg sysadmin.Config) bfs3.Prior {
	return func(task *rlglue.TaskSpec) (prior bayes.BeliefState) {
		mdp := sysadmin.NewSysMDP(cfg)
		transition := &bayes.MDPTransition{mdp}
		reward := &bayes.MDPReward{mdp}
		terminal := &bayes.MDPTerminal{transition}
		prior = bayes.NewBelief(0, reward, transition, terminal, nil)
		return
	}
}

type Config struct {
	Sysadmin	sysadmin.Config
	BFS3		bfs3.Config
}

func ConfigDefault() (cfg Config) {
	cfg.Sysadmin = sysadmin.ConfigDefault()
	cfg.BFS3 = bfs3.ConfigDefault()
	return
}
func main() {
	defer nicetrace.Print()
	cfg := ConfigDefault()
	argcfg.LoadArgs(&cfg)
	agent := sysadmin.NewBFS3Agent(cfg.BFS3, cfg.Sysadmin)
	rlglue.LoadAgent(agent)
}
