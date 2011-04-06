package sysadmin

import (
	"fmt"
	"go-glue.googlecode.com/hg/rlglue"
	"go-glue.googlecode.com/hg/rltools/discrete"
)

func NewMDPEnvironment(cfg Config) (env rlglue.Environment) {
	mdp := NewSysMDP(cfg)
	env = discrete.NewMDPEnvironment(mdp, mdp.Task, (1<<uint(cfg.NumSystems))-1)
	return
}

type SysMDP struct {
	Cfg			Config
	Task			*rlglue.TaskSpec
	maxStates, maxActions	uint64
	t			[][]float64
	r			[]float64
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
	for s := uint64(0); s < this.maxStates; s++ {
		for a := uint64(0); a < this.maxActions; a++ {
			k := s + a*this.maxStates
			this.t[k] = make([]float64, this.maxStates)
			this.r[k] = this.computeR(s, a)
			for n := uint64(0); n < this.maxStates; n++ {
				this.t[k][n] = this.computeT(s, a, n)
			}
		}
	}
	return
}
func (this *SysMDP) computeT(s, a, n uint64) (p float64) {
	sv := this.Task.Obs.Ints.Values(s)
	nv := this.Task.Obs.Ints.Values(n)
	p = 1
	for i, no := range nv {
		var fp float64
		if a == uint64(i) {
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
		if sv[i] == 1 || a == uint64(i) {
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
func (this *SysMDP) computeR(s, a uint64) (r float64) {
	sv := this.Task.Obs.Ints.Values(s)
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
func (this *SysMDP) S() uint64 {
	return this.maxStates
}
func (this *SysMDP) A() uint64 {
	return this.maxActions
}
func (this *SysMDP) T(s, a, n uint64) float64 {
	k := s + a*this.maxStates
	return this.t[k][n]
}
func (this *SysMDP) R(s, a uint64) float64 {
	k := s + a*this.maxStates
	return this.r[k]
}
func (this *SysMDP) GetGamma() float64 {
	return this.Cfg.DiscountFactor
}
