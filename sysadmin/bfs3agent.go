package sysadmin

import (
	"go-glue.googlecode.com/hg/rlglue"
	"go-glue.googlecode.com/hg/rlalg/bfs3"
	"github.com/skelterjohn/rlbayes"
	"go-glue.googlecode.com/hg/rltools/discrete"
)

func GetConfiguredPrio2r(cfg Config) bfs3.Prior {
	return func(task *rlglue.TaskSpec) bayes.BeliefState {
		cfg.NumSystems = len(task.Obs.Ints)
		env := New(cfg)
		env.EnvInit()
		env.EnvStart()
		return env
	}
}
func GetConfiguredPrior(cfg Config) bfs3.Prior {
	return func(task *rlglue.TaskSpec) (prior bayes.BeliefState) {
		mdp := NewSysMDP(cfg)
		transition := &bayes.MDPTransition{mdp}
		reward := &bayes.MDPReward{mdp}
		terminal := &bayes.MDPTerminal{transition}
		prior = bayes.NewBelief(0, reward, transition, terminal, nil)
		return
	}
}
func NewBFS3Agent(bfs3cfg bfs3.Config, cfg Config) (agent *bfs3.BFS3Agent) {
	prior := GetConfiguredPrior(cfg)
	agent = bfs3.New(prior)
	agent.Cfg = bfs3cfg
	return
}
func (this *Environment) Hashcode() (hash uint64) {
	hash = this.hash
	return
}
func (this *Environment) LessThan(oi interface{}) bool {
	other := oi.(*Environment)
	return this.hash < other.hash
}
func (this *Environment) Next(action uint64) (o discrete.Oracle, r float64) {
	act := rlglue.NewAction([]int32{int32(action)}, []float64{}, []byte{})
	next := new(Environment)
	*next = *this
	next.status = append([]bool{}, this.status...)
	_, r, _ = next.EnvStep(act)
	o = next
	return
}
func (this *Environment) Terminal() bool {
	return false
}
func (this *Environment) Update(action uint64, state uint64, reward float64) (nbs bayes.BeliefState) {
	next := new(Environment)
	*next = *this
	next.status = append([]bool{}, this.status...)
	next.Teleport(state)
	nbs = next
	return
}
func (this *Environment) UpdateTerminal(action uint64, reward float64) (next bayes.BeliefState) {
	return this
}
func (this *Environment) Teleport(state uint64) {
	ints := this.task.Obs.Ints.Values(state)
	for i, v := range ints {
		this.status[i] = v == 1
	}
	this.hash = state
}
func (this *Environment) GetState() uint64 {
	return this.hash
}
