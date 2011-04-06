package gbandits

import (
	"unsafe"
	"math"
	"fmt"
	"strings"
	"strconv"
	"rand"
	"gostat.googlecode.com/hg/stat"
	"go-glue.googlecode.com/hg/rlglue"
	"github.com/skelterjohn/rlbayes"
	"go-glue.googlecode.com/hg/rltools/discrete"
)

type Belief struct {
	Totals				[]float64
	Counts				[]float64
	mu0, sigmasqr0, sigmasqr1	float64
	M				uint64
	generators			[]func() float64
	hash				uint64
}

func NewBelief(numActions uint64, mu0, sigmasqr0, sigmasqr1 float64, M uint64) (this *Belief) {
	this = new(Belief)
	this.Totals = make([]float64, numActions)
	this.Counts = append([]float64{}, this.Totals...)
	this.mu0, this.sigmasqr0, this.sigmasqr1 = mu0, sigmasqr0, sigmasqr1
	this.M = M
	this.generators = make([]func() float64, numActions)
	this.rehash()
	return
}
func (this *Belief) rehash() {
	this.hash = 0
	for i, _ := range this.Counts {
		a := this.Counts[i]
		b := this.Totals[i]
		ai := *(*uint64)(unsafe.Pointer(&a))
		bi := *(*uint64)(unsafe.Pointer(&b))
		this.hash += ai << uint(i)
		this.hash += bi << uint(i+1)
	}
}
func (this *Belief) Hashcode() (hash uint64) {
	return this.hash
}
func (this *Belief) LessThan(other interface{}) bool {
	ob := other.(*Belief)
	for i, c := range this.Totals {
		if c < ob.Totals[i] {
			return true
		}
		if this.Counts[i] < ob.Counts[i] {
			return true
		}
	}
	return false
}
func (this *Belief) Next(action uint64) (o discrete.Oracle, r float64) {
	if this.generators[action] == nil {
		sigmasqrX := 1 / (1/this.sigmasqr0 + this.Counts[action]/this.sigmasqr1)
		muX := this.mu0/this.sigmasqr0 + this.Totals[action]/this.sigmasqr1
		muX *= sigmasqrX
		sigmaX := math.Sqrt(sigmasqrX)
		genMu := stat.Normal(muX, sigmaX)
		this.generators[action] = func() float64 {
			mu := genMu()
			return stat.NextNormal(mu, math.Sqrt(this.sigmasqr1))
		}
	}
	r = this.generators[action]()
	o = this.Update(action, 0, r)
	return
}
func (this *Belief) Teleport(state uint64) {
}
func (this *Belief) Terminal() bool {
	return false
}
func (this *Belief) GetState() uint64 {
	return 0
}
func (this *Belief) Update(action uint64, state uint64, reward float64) (nextb bayes.BeliefState) {
	if this.M != 0 && this.Counts[action] >= float64(this.M) {
		return this
	}
	next := new(Belief)
	*next = *this
	next.Totals = append([]float64{}, this.Totals...)
	next.Counts = append([]float64{}, this.Counts...)
	next.generators = make([]func() float64, len(this.Totals))
	next.Totals[action] += reward
	next.Counts[action] += 1
	next.rehash()
	nextb = next
	return
}
func (this *Belief) UpdateTerminal(action uint64, reward float64) (next bayes.BeliefState) {
	next = this
	return
}

type Env struct {
	numActions	int
	belief		*Belief
	obs		rlglue.Observation
}

func NewEnv(belief *Belief) (this *Env) {
	this = new(Env)
	this.numActions = len(belief.Counts)
	this.belief = belief
	this.obs = rlglue.NewObservation([]int32{0}, nil, nil)
	return
}
func (this *Env) EnvInit() string {
	format := "VERSION RL-Glue-3.0 PROBLEMTYPE episodic DISCOUNTFACTOR 0.95 OBSERVATIONS INTS (0 0) ACTIONS INTS (0 %d) REWARDS (-Inf Inf)"
	tstr := fmt.Sprintf(format, this.numActions-1)
	return tstr
}
func (this *Env) EnvStart() (obs rlglue.Observation) {
	return this.obs
}
func (this *Env) EnvStep(action rlglue.Action) (obs rlglue.Observation, r float64, t bool) {
	obs = this.obs
	var o discrete.Oracle
	a := uint64(action.Ints()[0])
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
