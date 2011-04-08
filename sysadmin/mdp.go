package sysadmin

import (
	"fmt"
	"go-glue.googlecode.com/hg/rlglue"
	"go-glue.googlecode.com/hg/rltools"
	"go-glue.googlecode.com/hg/rltools/discrete"
)

func NewMDPEnvironment(cfg Config) (env rlglue.Environment) {
	mdp := NewSysMDP(cfg)
	env = discrete.NewMDPEnvironment(mdp, mdp.Task, (1<<uint(cfg.NumSystems))-1)
	return
}

type SysMDP struct {
	Cfg                   Config
	Task                  *rlglue.TaskSpec
	maxStates, maxActions uint64
	t                     [][]float64
	r                     []float64
}

func NewSysMDP(Cfg Config) (this *SysMDP) {
	this = &SysMDP{Cfg: Cfg}
	fstr := "VERSION RL-Glue-3.0 PROBLEMTYPE episodic DISCOUNTFACTOR %f OBSERVATIONS INTS (%d 0 1) ACTIONS INTS (0 %d) REWARDS (-1.0 1.0)"
	taskString := fmt.Sprintf(fstr, this.Cfg.DiscountFactor, this.Cfg.NumSystems, this.Cfg.NumSystems)
	this.Task, _ = rlglue.ParseTaskSpec(taskString)
	this.maxStates = this.Task.Obs.Ints.Count()
	this.maxActions = this.Task.Act.Ints.Count()
	this.t = make([][]float64, this.maxStates*this.maxActions)
	this.r = make([]float64, this.maxStates*this.maxActions)
	for s := range this.S64() {
		for a := range this.A64() {
			k := s.Hashcode() + a.Hashcode()*this.maxStates
			this.t[k] = make([]float64, this.maxStates)
			this.r[k] = this.computeR(s, a)
			for n := range this.S64() {
				this.t[k][n] = this.computeT(s, a, n)
			}
		}
	}
	return
}
func (this *SysMDP) computeT(s discrete.State, a discrete.Action, n discrete.State) (p float64) {
	sv := this.Task.Obs.Ints.Values(s.Hashcode())
	nv := this.Task.Obs.Ints.Values(n.Hashcode())
	p = 1
	for i, no := range nv {
		var fp float64
		if a == discrete.Action(i) {
			fp = 0
		} else {
			fp = this.Cfg.FailBase
			li := (i + this.Cfg.NumSystems - 1) % this.Cfg.NumSystems
			ri := (i + 1) % this.Cfg.NumSystems
			ls := sv[li] == 1
			rs := sv[ri] == 1
			if li < i {
				ls = nv[li] == 1
			}
			if !ls {
				fp += this.Cfg.FailIncr
			}
			if !rs {
				fp += this.Cfg.FailIncr
			}
		}
		if sv[i] == 1 || a == discrete.Action(i) {
			if no == 0 {
				p *= fp
			} else {
				p *= 1 - fp
			}
		} else {
			if no == 0 {
				p *= this.Cfg.FailStay
			} else {
				p *= 1 - this.Cfg.FailStay
			}
		}
	}
	return
}
func (this *SysMDP) computeR(s discrete.State, a discrete.Action) (r float64) {
	sv := this.Task.Obs.Ints.Values(s.Hashcode())
	for _, v := range sv {
		r += float64(v)
	}
	if int(a) < this.Cfg.NumSystems {
		r -= 1
	}
	return
}
func (this *SysMDP) GetTask() (task *rlglue.TaskSpec) {
	task = this.Task
	return
}
func (this *SysMDP) NumStates() uint64 {
	return this.maxStates
}
func (this *SysMDP) NumActions() uint64 {
	return this.maxActions
}
func (this *SysMDP) S() <-chan rltools.State {
	return discrete.AllStates(this.NumStates())
}
func (this *SysMDP) A() <-chan rltools.Action {
	return discrete.AllActions(this.NumActions())
}
func (this *SysMDP) S64() <-chan discrete.State {
	return discrete.AllStates64(this.NumStates())
}
func (this *SysMDP) A64() <-chan discrete.Action {
	return discrete.AllActions64(this.NumActions())
}
func (this *SysMDP) T(s discrete.State, a discrete.Action, n discrete.State) float64 {
	k := s.Hashcode() + a.Hashcode()*this.maxStates
	return this.t[k][n]
}
func (this *SysMDP) R(s discrete.State, a discrete.Action) float64 {
	k := s.Hashcode() + a.Hashcode()*this.maxStates
	return this.r[k]
}
func (this *SysMDP) GetGamma() float64 {
	return this.Cfg.DiscountFactor
}
