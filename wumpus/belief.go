package wumpus

import (
	"fmt"
	"go-glue.googlecode.com/hg/rlglue"
	"github.com/skelterjohn/rlbayes"
	"go-glue.googlecode.com/hg/rltools/discrete"
)

var taskstr = fmt.Sprintf("VERSION RL-Glue-3.0 PROBLEMTYPE episodic DISCOUNTFACTOR %f OBSERVATIONS INTS (2 0 3) (0 3) (16 0 %d) ACTIONS INTS (0 3) REWARDS (%f, 1)", Gamma, len(SquareMapToValue), -.01/(1-Gamma))
var task, _ = rlglue.ParseTaskSpec(taskstr)

type Belief struct {
	Hunter
	term bool
	hash uint64
}

func NewBelief(mb MapBelief) (this *Belief) {
	this = new(Belief)
	this.x = 0
	this.y = 0
	this.dir = 1
	this.belief = mb
	return
}
func (this *Belief) Hashcode() uint64 {
	return this.hash
}
func (this *Belief) LessThan(oi interface{}) bool {
	return this.hash < oi.(*Belief).hash
}
func (this *Belief) Next(action discrete.Action) (o discrete.Oracle, r float64) {
	var nexthunter Hunter
	var t bool
	switch action {
	case 0:
		nexthunter, r, t = this.TurnLeft()
	case 1:
		nexthunter, r, t = this.TurnRight()
	case 2:
		nexthunter, r, t = this.Move()
	case 3:
		nexthunter, r, t = this.Shoot()
	}
	if t {
		ob := &Belief{}
		ob.term = true
		o = ob
		return
	}
	ob := &Belief{Hunter: nexthunter, term: false}
	ob.hash = ob.GetState().Hashcode()
	o = ob
	return
}
func (this *Belief) Terminal() bool {
	return this.term
}
func (this *Belief) Update(action discrete.Action, state discrete.State, reward float64) (next bayes.BeliefState) {
	nb := new(Belief)
	*nb = *this
	nb.belief = append(MapBelief{}, this.belief...)
	nb.Teleport(state)
	next = nb
	return
}
func (this *Belief) UpdateTerminal(action discrete.Action, reward float64) (next bayes.BeliefState) {
	return this
}
func (this *Belief) Teleport(state discrete.State) {
	this.hash = state.Hashcode()
	indices := task.Obs.Ints.Values(state.Hashcode())
	for i, ii := range indices[3:] {
		indices[i+3] = GetValue(ii)
	}
	this.x = indices[0]
	this.y = indices[1]
	this.dir = indices[2]
	this.belief = indices[3:]
	return
}
func (this *Belief) GetState() discrete.State {
	values := make([]int32, len(task.Obs.Ints))
	values[0] = this.x
	values[1] = this.y
	values[2] = this.dir
	for i, v := range this.belief {
		values[i+3] = GetIndex(v)
	}
	return discrete.State(task.Obs.Ints.Index(values))
}
