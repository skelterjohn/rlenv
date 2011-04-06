package paint

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"gostat.googlecode.com/hg/stat"
	"go-glue.googlecode.com/hg/rlglue"
	"go-glue.googlecode.com/hg/rltools/discrete"
	"github.com/skelterjohn/rlbayes"
)

type Config struct{ NumCans int }

func ConfigDefault() (cfg Config) {
	cfg.NumCans = 1
	return
}

type Can struct{ Painted, Polished, Scratched, Done bool }
type Env struct {
	Cans	[]Can
	Live	bool
}

func New(cfg Config) (this *Env) {
	this = new(Env)
	this.Cans = make([]Can, cfg.NumCans)
	return
}
func (this *Env) Log(format string, params ...interface{}) {
	if this.Live {
		fmt.Fprintf(os.Stderr, format, params...)
	}
}
func (this *Env) EnvInit() (taskstring string) {
	format := "VERSION RL-Glue-3.0 PROBLEMTYPE episodic DISCOUNTFACTOR 0.95 OBSERVATIONS INTS (%d 0 1) ACTIONS INTS (0 %d) (0 3) REWARDS (-1.0 10.0)"
	taskstring = fmt.Sprintf(format, 4*len(this.Cans), len(this.Cans)-1)
	return
}
func (this *Env) makeObs() (obs rlglue.Observation) {
	ints := make([]int32, 4*len(this.Cans))
	for i, can := range this.Cans {
		if can.Painted {
			ints[i*4] = 1
		}
		if can.Polished {
			ints[i*4+1] = 1
		}
		if can.Scratched {
			ints[i*4+2] = 1
		}
		if can.Done {
			ints[i*4+3] = 1
		}
	}
	obs = rlglue.NewObservation(ints, []float64{}, []byte{})
	return
}
func (this *Env) EnvStart() (obs rlglue.Observation) {
	for i := range this.Cans {
		this.Cans[i] = Can{}
	}
	obs = this.makeObs()
	return
}
func (this *Env) EnvStep(action rlglue.Action) (obs rlglue.Observation, r float64, t bool) {
	whichCan := action.Ints()[0]
	process := action.Ints()[1]
	t = true
	for _, can := range this.Cans {
		if !can.Done {
			t = false
		}
	}
	this.Log("%v ", this.Cans)
	if t {
		this.Log("finished\n\n")
		obs = this.makeObs()
		r = 0
		return
	}
	can := this.Cans[whichCan]
	r = -1
	if !can.Done {
		switch process {
		case 0:
			this.Log("painting can %d\n", whichCan+1)
			outcome := stat.NextChoice([]float64{.6, .3, .1})
			switch outcome {
			case 0:
				can.Painted = true
			case 1:
				can.Painted = true
				can.Scratched = true
			case 2:
			}
		case 1:
			this.Log("polishing can %d\n", whichCan+1)
			outcome := stat.NextChoice([]float64{.2, .2, .3, .2, .1})
			switch outcome {
			case 0:
				can.Painted = false
			case 1:
				can.Scratched = false
			case 2:
				can.Polished = true
				can.Painted = false
				can.Scratched = false
			case 3:
				can.Polished = true
				can.Painted = false
			case 4:
			}
		case 2:
			this.Log("shortcut can %d\n", whichCan+1)
			outcome := stat.NextChoice([]float64{0.05, 0.95})
			switch outcome {
			case 0:
				can.Painted = true
				can.Polished = true
			case 1:
			}
		case 3:
			this.Log("finishing can %d\n", whichCan+1)
			if can.Painted && can.Polished && !can.Scratched && !can.Done {
				can.Done = true
				r = 10
			} else {
				t = true
				r = -100000
			}
		}
		this.Cans[whichCan] = can
	}
	obs = this.makeObs()
	return
}
func (this *Env) ActionAvailable(action uint64) bool {
	which := int(action) % len(this.Cans)
	how := int(action) / len(this.Cans)
	if how == 3 {
		if !this.Cans[which].Painted {
			return false
		}
		if !this.Cans[which].Polished {
			return false
		}
		if this.Cans[which].Scratched {
			return false
		}
		if this.Cans[which].Done {
			return false
		}
	}
	return true
}
func (this *Env) EnvCleanup() {
}
func (this *Env) EnvMessage(message string) (reply string) {
	tokens := strings.Split(message, " ", -1)
	if tokens[0] == "seed" {
		seed, _ := strconv.Atoi64(tokens[1])
		stat.Seed(seed)
	}
	return ""
}

type Oracle struct {
	Env
	isTerminal	bool
	Task		*rlglue.TaskSpec
	hash		uint64
}

func NewOracle(env Env) (this *Oracle) {
	this = new(Oracle)
	this.Env = env
	this.Task, _ = rlglue.ParseTaskSpec(this.Env.EnvInit())
	this.rehash()
	return
}
func (this *Oracle) rehash() {
	this.hash = this.Task.Obs.Ints.Index(this.Env.makeObs().Ints())
}
func (this *Oracle) Hashcode() (hash uint64) {
	return this.hash
}
func (this *Oracle) Equals(other interface{}) bool {
	oo := other.(*Oracle)
	return this.Hashcode() == oo.Hashcode()
}
func (this *Oracle) LessThan(other interface{}) bool {
	oo := other.(*Oracle)
	return this.Hashcode() < oo.Hashcode()
}
func (this *Oracle) Next(action uint64) (o discrete.Oracle, r float64) {
	avalues := this.Task.Act.Ints.Values(action)
	act := rlglue.NewAction(avalues, []float64{}, []byte{})
	next := new(Oracle)
	*next = *this
	next.Cans = append([]Can{}, this.Cans...)
	_, r, next.isTerminal = next.Env.EnvStep(act)
	next.rehash()
	o = next
	return
}
func (this *Oracle) Terminal() bool {
	return this.isTerminal
}
func (this *Oracle) Update(action, state uint64, reward float64) (next bayes.BeliefState) {
	no := new(Oracle)
	*no = *this
	no.Cans = append([]Can{}, this.Cans...)
	no.Teleport(state)
	next = no
	return
}
func (this *Oracle) UpdateTerminal(action uint64, reward float64) (next bayes.BeliefState) {
	return this
}
func (this *Oracle) Teleport(state uint64) {
	ints := this.Task.Obs.Ints.Values(state)
	for i := range this.Cans {
		this.Cans[i].Painted = ints[i*4] == 1
		this.Cans[i].Polished = ints[i*4+1] == 1
		this.Cans[i].Scratched = ints[i*4+2] == 1
		this.Cans[i].Done = ints[i*4+3] == 1
	}
	this.hash = state
}
func (this *Oracle) GetState() (state uint64) {
	return this.hash
}
