package bbandits

import (
	"strings"
	"strconv"
	"rand"
	"unsafe"
	"gostat.googlecode.com/hg/stat"
	"fmt"
	"github.com/skelterjohn/rlbayes"
	"go-glue.googlecode.com/hg/rlglue"
	"go-glue.googlecode.com/hg/rltools/discrete"
)

type Belief struct {
	counts      []float64
	totals      []float64
	visitOffset []float64
	M           uint64
	hash        uint64
}

func NewBelief(alphas, betas []float64) (this *Belief) {
	this = new(Belief)
	this.counts = make([]float64, len(alphas))
	this.totals = make([]float64, len(alphas))
	copy(this.counts, alphas)
	for i, alpha := range alphas {
		this.totals[i] = alpha + betas[i]
	}
	this.visitOffset = append([]float64{}, this.totals...)
	this.rehash()
	return
}
func (this *Belief) rehash() {
	this.hash = 0
	for i, _ := range this.counts {
		a := this.counts[i]
		b := this.totals[i] - a
		ai := *(*uint64)(unsafe.Pointer(&a))
		bi := *(*uint64)(unsafe.Pointer(&b))
		this.hash += ai << uint(i)
		this.hash += bi << uint(i)
	}
}
func (this *Belief) Hashcode() (hash uint64) {
	return this.hash
}
func (this *Belief) LessThan(other interface{}) bool {
	ob := other.(*Belief)
	for i, c := range this.totals {
		if c < ob.totals[i] {
			return true
		}
		if this.counts[i] < ob.counts[i] {
			return true
		}
	}
	return false
}
func (this *Belief) Next(c discrete.Action) (o discrete.Oracle, r float64) {
	visits := this.totals[c] - this.visitOffset[c]
	if this.M != 0 && uint64(visits) >= this.M {
		o = this
		if this.counts[c]/this.totals[c] > stat.NextUniform() {
			r = 1
		}
		return
	}
	next := new(Belief)
	next.counts = make([]float64, len(this.counts))
	copy(next.counts, this.counts)
	next.totals = make([]float64, len(this.totals))
	copy(next.totals, this.totals)
	next.visitOffset = this.visitOffset
	prob := this.counts[c] / this.totals[c]
	if prob > stat.NextUniform() {
		r = 1
		next.counts[c] += 1
	}
	next.totals[c] += 1
	next.rehash()
	next.M = this.M
	o = next
	return
}
func (this *Belief) Teleport(state discrete.State) {
}
func (this *Belief) Terminal() bool {
	return false
}
func (this *Belief) GetState() discrete.State {
	return 0
}
func (this *Belief) Update(action discrete.Action, state discrete.State, reward float64) (nextb bayes.BeliefState) {
	visits := this.totals[action] - this.visitOffset[action]
	if this.M != 0 && uint64(visits) >= this.M {
		return this
	}
	next := new(Belief)
	next.counts = make([]float64, len(this.counts))
	copy(next.counts, this.counts)
	next.totals = make([]float64, len(this.totals))
	copy(next.totals, this.totals)
	next.visitOffset = this.visitOffset
	if reward == 1 {
		next.counts[action] += 1
	}
	next.totals[action] += 1
	next.M = this.M
	next.rehash()
	nextb = next
	return
}
func (this *Belief) UpdateTerminal(action discrete.Action, reward float64) (next bayes.BeliefState) {
	next = this
	return
}

type Env struct {
	numActions int
	belief     *Belief
	obs        rlglue.Observation
}

func NewEnv(belief *Belief) (this *Env) {
	this = new(Env)
	this.numActions = len(belief.counts)
	this.belief = belief
	this.obs = rlglue.NewObservation([]int32{0}, nil, nil)
	return
}
func (this *Env) EnvInit() string {
	format := "VERSION RL-Glue-3.0 PROBLEMTYPE episodic DISCOUNTFACTOR 0.95 OBSERVATIONS INTS (0 0) ACTIONS INTS (0 %d) REWARDS (0 1)"
	tstr := fmt.Sprintf(format, this.numActions-1)
	return tstr
}
func (this *Env) EnvStart() (obs rlglue.Observation) {
	obs = this.obs
	return
}
func (this *Env) EnvStep(action rlglue.Action) (obs rlglue.Observation, r float64, t bool) {
	obs = this.obs
	var o discrete.Oracle
	a := discrete.Action(action.Ints()[0])
	o, r = this.belief.Next(a)
	this.belief = o.(*Belief)
	t = false
	return
}
func (this *Env) EnvCleanup() {
}
func (this *Env) EnvMessage(message string) (reply string) {
	tokens := strings.Split(message, " ", -1)
	if tokens[0] == "seed" {
		seed, _ := strconv.Atoi64(tokens[1])
		rand.Seed(seed)
	}
	return ""
}
