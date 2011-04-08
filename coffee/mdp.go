package coffee

import (
	"go-glue.googlecode.com/hg/rlglue"
	"go-glue.googlecode.com/hg/rltools"
	"go-glue.googlecode.com/hg/rltools/discrete"
	"os"
	"fmt"
	"math"
)

const (
	IndexInOffice = iota
	IndexWet
	IndexUmbrella
	IndexRaining
	IndexRobotHasCoffee
	IndexUserHasCoffee
)

type MDP struct {
	Cfg                   Config
	Task                  *rlglue.TaskSpec
	maxStates, maxActions uint64
	svs                   [][]int32
	t                     [][]float64
	r                     []float64
}

func NewMDP(cfg Config) (this *MDP) {
	this = new(MDP)
	this.Cfg = cfg
	fstr := "VERSION RL-Glue-3.0 PROBLEMTYPE episodic DISCOUNTFACTOR %f OBSERVATIONS INTS (6 0 1) ACTIONS INTS (0 3) REWARDS (%f %f)"
	rmin := 0.0
	rmax := 1.0
	taskString := fmt.Sprintf(fstr, this.Cfg.DiscountFactor, rmin, rmax)
	this.Task, _ = rlglue.ParseTaskSpec(taskString)
	this.maxStates = this.Task.Obs.Ints.Count()
	this.maxActions = this.Task.Act.Ints.Count()
	this.svs = make([][]int32, this.maxStates)
	for s := uint64(0); s < this.maxStates; s++ {
		this.svs[s] = this.Task.Obs.Ints.Values(s)
	}
	this.t = make([][]float64, this.maxStates*this.maxActions)
	this.r = make([]float64, this.maxStates*this.maxActions)
	for s := range this.S64() {
		for a := range this.A64() {
			k := s.Hashcode() + this.maxStates*a.Hashcode()
			this.r[k] = this.computeR(s, a)
			this.t[k] = make([]float64, this.maxStates)
			sum := 0.0
			for n := range this.S64() {
				this.t[k][n] = this.computeT(s, a, n)
				sum += this.t[k][n]
			}
			if sum > 1 {
				fmt.Fprintf(os.Stderr, "%v %v\n", this.svs[s], a)
				for n := range this.S64() {
					if this.t[k][n] != 0 {
						fmt.Fprintf(os.Stderr, "\t%v : %v\n", this.svs[n], this.t[k][n])
					}
				}
			}
		}
	}
	return
}
func (this *MDP) computeR(s discrete.State, a discrete.Action) (r float64) {
	sv := this.svs[s]
	Wet := sv[IndexWet] == 1
	UserHasCoffee := sv[IndexUserHasCoffee] == 1
	if !Wet {
		r += this.Cfg.RewardDry
	}
	if UserHasCoffee {
		r += this.Cfg.RewardUserHasCoffee
	}
	return
}
func (this *MDP) computeT(s discrete.State, a discrete.Action, n discrete.State) (p float64) {
	sv := this.svs[s]
	nv := this.svs[n]
	same := func(f bool) (p float64) {
		if f {
			p = 1
		} else {
			p = 0
		}
		return
	}
	_ = same
	InOffice := sv[IndexInOffice] == 1
	Wet := sv[IndexWet] == 1
	Raining := sv[IndexRaining] == 1
	Umbrella := sv[IndexUmbrella] == 1
	RobotHasCoffee := sv[IndexRobotHasCoffee] == 1
	UserHasCoffee := sv[IndexUserHasCoffee] == 1
	if UserHasCoffee {
		if nv[IndexUserHasCoffee] == 1 {
			return 0
		}
		return math.Pow(.5, 5)
	}
	var pUserHasCoffee, pRobotHasCoffee, pWet, pRaining, pUmbrella, pInOffice float64
	switch a {
	case Move:
		pUserHasCoffee = same(UserHasCoffee)
		pRobotHasCoffee = same(RobotHasCoffee)
		if Wet {
			pWet = 1
		} else {
			if Raining {
				if Umbrella {
					pWet = this.Cfg.ProbGetWetUmbrella
				} else {
					pWet = this.Cfg.ProbGetWetNoUmbrella
				}
			} else {
				pWet = 0
			}
		}
		pRaining = same(Raining)
		pUmbrella = same(Umbrella)
		if InOffice {
			pInOffice = 1 - this.Cfg.ProbMove
		} else {
			pInOffice = this.Cfg.ProbMove
		}
	case DeliverCoffee:
		if UserHasCoffee {
			pUserHasCoffee = 1
		} else {
			if RobotHasCoffee {
				if InOffice {
					pUserHasCoffee = this.Cfg.ProbDeliverCoffee
				} else {
					pUserHasCoffee = 0
				}
			} else {
				pUserHasCoffee = 0
			}
		}
		if RobotHasCoffee {
			if InOffice {
				pRobotHasCoffee = this.Cfg.ProbKeepCoffeeOffice
			} else {
				pRobotHasCoffee = this.Cfg.ProbKeepCoffeeNoOffice
			}
		} else {
			pRobotHasCoffee = 0
		}
		pWet = same(Wet)
		pRaining = same(Raining)
		pUmbrella = same(Umbrella)
		pInOffice = same(InOffice)
	case GetUmbrella:
		pUserHasCoffee = same(UserHasCoffee)
		pRobotHasCoffee = same(RobotHasCoffee)
		pWet = same(Wet)
		pRaining = same(Raining)
		if Umbrella {
			pUmbrella = 1
		} else {
			if InOffice {
				pUmbrella = this.Cfg.ProbFindUmbrella
			} else {
				pUmbrella = 0
			}
		}
		pInOffice = same(InOffice)
	case BuyCoffee:
		pUserHasCoffee = same(UserHasCoffee)
		if RobotHasCoffee {
			pRobotHasCoffee = 1
		} else {
			if InOffice {
				pRobotHasCoffee = 0
			} else {
				pRobotHasCoffee = this.Cfg.ProbBuyCoffee
			}
		}
		pWet = same(Wet)
		pRaining = same(Raining)
		pUmbrella = same(Umbrella)
		pInOffice = same(InOffice)
	}
	flip := func(f bool, p float64) (rp float64) {
		rp = p
		if !f {
			rp = 1 - p
		}
		return
	}
	InOffice = nv[IndexInOffice] == 1
	Wet = nv[IndexWet] == 1
	Raining = nv[IndexRaining] == 1
	Umbrella = nv[IndexUmbrella] == 1
	RobotHasCoffee = nv[IndexRobotHasCoffee] == 1
	UserHasCoffee = nv[IndexUserHasCoffee] == 1
	p = 1
	p *= flip(UserHasCoffee, pUserHasCoffee)
	p *= flip(RobotHasCoffee, pRobotHasCoffee)
	p *= flip(InOffice, pInOffice)
	p *= flip(Wet, pWet)
	p *= flip(Raining, pRaining)
	p *= flip(Umbrella, pUmbrella)
	return
}
func (this *MDP) GetTask() (task *rlglue.TaskSpec) {
	task = this.Task
	return
}
func (this *MDP) NumStates() uint64 {
	return this.maxStates
}
func (this *MDP) NumActions() uint64 {
	return this.maxActions
}
func (this *MDP) S() <-chan rltools.State {
	return discrete.AllStates(this.maxStates)
}
func (this *MDP) A() <-chan rltools.Action {
	return discrete.AllActions(this.maxActions)
}
func (this *MDP) S64() <-chan discrete.State {
	return discrete.AllStates64(this.maxStates)
}
func (this *MDP) A64() <-chan discrete.Action {
	return discrete.AllActions64(this.maxActions)
}
func (this *MDP) T(s discrete.State, a discrete.Action, n discrete.State) float64 {
	k := s.Hashcode() + a.Hashcode()*this.maxStates
	return this.t[k][n]
}
func (this *MDP) R(s discrete.State, a discrete.Action) float64 {
	k := s.Hashcode() + a.Hashcode()*this.maxStates
	return this.r[k]
}
func (this *MDP) GetGamma() float64 {
	return this.Cfg.DiscountFactor
}
