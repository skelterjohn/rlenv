package main

import (
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"go-glue.googlecode.com/hg/rlglue"
	"github.com/skelterjohn/rlbayes"
	"github.com/skelterjohn/rlalg/bfs3"
	"github.com/skelterjohn/rlenv/coffee"
)

func GetTruthFunc(cfg coffee.Config) bfs3.Prior {
	return func(task *rlglue.TaskSpec) (prior bayes.BeliefState) {
		mdp := coffee.NewMDP(cfg)
		transition := &bayes.MDPTransition{mdp}
		reward := &bayes.MDPReward{mdp}
		terminal := &bayes.MDPTerminal{transition}
		prior = bayes.NewBelief(0, reward, transition, terminal, nil)
		return
	}
}

type Config struct {
	Coffee	coffee.Config
	BFS3	bfs3.Config
}

func ConfigDefault() (cfg Config) {
	cfg.Coffee = coffee.ConfigDefault()
	cfg.BFS3 = bfs3.ConfigDefault()
	return
}
func NewBFS3Agent(cfg Config) (agent *bfs3.BFS3Agent) {
	prior := GetTruthFunc(cfg.Coffee)
	agent = bfs3.New(prior)
	agent.Cfg = cfg.BFS3
	return
}
func main() {
	defer nicetrace.Print()
	cfg := ConfigDefault()
	argcfg.LoadArgs(&cfg)
	agent := NewBFS3Agent(cfg)
	rlglue.LoadAgent(agent)
}
