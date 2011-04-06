package wumpus

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"gostat.googlecode.com/hg/stat"
	"go-glue.googlecode.com/hg/rlglue"
)

type Env struct {
	hunter		Hunter
	observed	MapBelief
}

func New() (this *Env) {
	this = new(Env)
	return
}
func (this *Env) MakeObs() (obs rlglue.Observation) {
	k := this.hunter.x + this.hunter.y*4
	this.observed[k] = this.hunter.belief[k]
	intarray := []int32{this.hunter.x, this.hunter.y, this.hunter.dir}
	indices := make([]int32, 16)
	for i, v := range this.observed {
		indices[i] = GetIndex(v)
	}
	intarray = append(intarray, indices...)
	obs = rlglue.NewObservation(intarray, []float64{}, []byte{})
	return
}
func (this *Env) EnvInit() (taskString string) {
	return taskstr
}
func (this *Env) EnvStart() (obs rlglue.Observation) {
	this.observed = MapBelief{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown}
	this.observed[0] = Empty
	start := MapBelief{Empty, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown}
	this.hunter.belief, _ = start.SampleWorld()
	this.hunter.x = 0
	this.hunter.y = 0
	this.hunter.dir = 1
	obs = this.MakeObs()
	fmt.Fprintf(os.Stderr, "Sending back\n%v\n%v\n", this.hunter, this.observed)
	return
}
func (this *Env) EnvStep(action rlglue.Action) (obs rlglue.Observation, r float64, t bool) {
	lastObs := this.MakeObs()
	println(action.Ints()[0])
	switch action.Ints()[0] {
	case 0:
		this.hunter, r, t = this.hunter.TurnLeft()
	case 1:
		this.hunter, r, t = this.hunter.TurnRight()
	case 2:
		this.hunter, r, t = this.hunter.Move()
	case 3:
		this.hunter, r, t = this.hunter.Shoot()
	}
	if !t {
		obs = this.MakeObs()
		fmt.Fprintf(os.Stderr, "Sending back\n%v\n%v\n", this.hunter, this.observed)
	} else {
		obs = lastObs
	}
	return
}
func (this *Env) EnvCleanup() {
	return
}
func (this *Env) EnvMessage(message string) (reply string) {
	if strings.HasPrefix(message, "seed ") {
		seed, _ := strconv.Atoi(message[5:])
		stat.Seed(int64(seed))
	}
	return
}
