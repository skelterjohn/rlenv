package sysadmin

import (
	"fmt"
	"strings"
	"strconv"
	"go-glue.googlecode.com/hg/rlglue"
	"gostat.googlecode.com/hg/stat"
)

var count int

type Config struct {
	FailBase, FailIncr	float64
	FailStay, StartBoot	float64
	NumSystems		int
	DiscountFactor		float64
}

func ConfigDefault() (cfg Config) {
	cfg.FailBase = 0.05
	cfg.FailIncr = 0.3
	cfg.FailStay = 0.95
	cfg.StartBoot = 0.9
	cfg.NumSystems = 8
	cfg.DiscountFactor = 0.9
	return
}

type Environment struct {
	task	*rlglue.TaskSpec
	cfg	Config
	status	[]bool
	hash	uint64
}

func New(cfg Config) (this *Environment) {
	this = new(Environment)
	this.cfg = cfg
	return
}
func (this *Environment) ConstructObs() (obs rlglue.Observation) {
	ints := make([]int32, this.cfg.NumSystems)
	for i := range ints {
		if this.status[i] {
			ints[i] = 1
		}
	}
	obs = rlglue.NewObservation(ints, []float64{}, []byte{})
	this.hash = this.task.Obs.Ints.Index(ints)
	count++
	return
}
func (this *Environment) EnvInit() (taskString string) {
	fstr := "VERSION RL-Glue-3.0 PROBLEMTYPE episodic DISCOUNTFACTOR %f OBSERVATIONS INTS (%d 0 1) ACTIONS INTS (0 %d) REWARDS (-1.0 1.0)"
	taskString = fmt.Sprintf(fstr, this.cfg.DiscountFactor, this.cfg.NumSystems, this.cfg.NumSystems)
	this.task, _ = rlglue.ParseTaskSpec(taskString)
	this.status = make([]bool, this.cfg.NumSystems)
	for i := range this.status {
		this.status[i] = stat.NextBernoulli(this.cfg.StartBoot) == 1
	}
	return
}
func (this *Environment) EnvStart() (obs rlglue.Observation) {
	return this.ConstructObs()
}
func (this *Environment) EnvStep(action rlglue.Action) (obs rlglue.Observation, r float64, t bool) {
	fps := make([]float64, len(this.status))
	reboot := int(action.Ints()[0])
	for i := range this.status {
		if reboot == i {
			fps[i] = 0
		} else {
			fps[i] = this.cfg.FailBase
			li := (i + this.cfg.NumSystems - 1) % this.cfg.NumSystems
			ri := (i + 1) % this.cfg.NumSystems
			if !this.status[li] {
				fps[i] += this.cfg.FailIncr
			}
			if !this.status[ri] {
				fps[i] += this.cfg.FailIncr
			}
		}
		if this.status[i] || reboot == i {
			this.status[i] = stat.NextUniform() < (1 - fps[i])
		} else {
			this.status[i] = stat.NextUniform() < (1 - this.cfg.FailStay)
		}
		if this.status[i] {
			r++
		}
	}
	if reboot < this.cfg.NumSystems {
		r--
	}
	obs = this.ConstructObs()
	return
}
func (this *Environment) EnvCleanup() {
}
func (this *Environment) EnvMessage(message string) (reply string) {
	if strings.HasPrefix(message, "seed ") {
		seed, _ := strconv.Atoi(message[5:])
		stat.Seed(int64(seed))
	}
	return ""
}
