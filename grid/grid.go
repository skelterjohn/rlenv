package grid

import (
	"strings"
	"strconv"
	"rand"
	"fmt"
	"go-glue.googlecode.com/hg/rlglue"
)

type Cell struct{ X, Y int32 }

func (c Cell) Equals(o Cell) bool {
	return c.X == o.X && c.Y == o.Y
}

type Env struct {
	Pos           Cell
	Start         Cell
	Goal          Cell
	Width, Height int32
}

func New(width, height int32) (ge *Env) {
	ge = new(Env)
	ge.Width, ge.Height = width, height
	ge.Start = Cell{0, 0}
	ge.Goal = Cell{width - 1, height - 1}
	return
}
func (ge *Env) EnvInit() string {
	format := "VERSION RL-Glue-3.0 PROBLEMTYPE episodic DISCOUNTFACTOR 1 OBSERVATIONS INTS (0 %d) (0 %d) ACTIONS INTS (0 3) REWARDS (-1.0 0.0)"
	tstr := fmt.Sprintf(format, ge.Width-1, ge.Height-1)
	return tstr
}
func (ge *Env) EnvStart() (obs rlglue.Observation) {
	ge.Pos = ge.Start
	obs = rlglue.NewObservation([]int32{ge.Pos.X, ge.Pos.Y}, []float64{}, []byte{})
	return
}
func (ge *Env) EnvStep(action rlglue.Action) (obs rlglue.Observation, r float64, t bool) {
	t = ge.Pos.Equals(ge.Goal)
	r = -1
	if t {
		r = 0
	}
	dir := action.Ints()[0]
	u := rand.Float64()
	if u < .1 {
		dir += 1
	} else if u < .2 {
		dir += 3
	}
	dir %= 4
	newPos := ge.Pos
	switch dir {
	case 0:
		newPos.Y++
	case 1:
		newPos.X++
	case 2:
		newPos.Y--
	case 3:
		newPos.X--
	}
	if newPos.X < 0 {
		newPos.X = 0
	}
	if newPos.Y < 0 {
		newPos.Y = 0
	}
	if newPos.X >= ge.Width {
		newPos.X = ge.Width - 1
	}
	if newPos.Y >= ge.Height {
		newPos.Y = ge.Height - 1
	}
	ge.Pos = newPos
	obs = rlglue.NewObservation([]int32{ge.Pos.X, ge.Pos.Y}, []float64{}, []byte{})
	return
}
func (ge *Env) EnvCleanup() {
}
func (ge *Env) EnvMessage(message string) (reply string) {
	tokens := strings.Split(message, " ", -1)
	if tokens[0] == "seed" {
		seed, _ := strconv.Atoi64(tokens[1])
		rand.Seed(seed)
	}
	return ""
}
