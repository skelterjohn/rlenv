package wumpus

import (
	"fmt"
)

const Gamma = .9

type Hunter struct {
	belief MapBelief
	x, y   int32
	dir    int32
	count  int
}

func (this Hunter) TurnLeft() (next Hunter, r float64, t bool) {
	next = this
	next.count++
	next.dir = (next.dir + 3) % 4
	r = -.01
	return
}
func (this Hunter) TurnRight() (next Hunter, r float64, t bool) {
	next = this
	next.count++
	next.dir = (next.dir + 1) % 4
	r = -.01
	return
}
func (this Hunter) Move() (next Hunter, r float64, t bool) {
	next = this

	sampledTruth, _ := this.belief.SampleWorld()
	switch next.dir {
	case 0:
		if next.y > 0 {
			next.y--
		}
	case 1:
		if next.x < 3 {
			next.x++
		}
	case 2:
		if next.y < 3 {
			next.y++
		}
	case 3:
		if next.x > 0 {
			next.x--
		}
	}

	nk := next.x + next.y*4
	next.belief = append(MapBelief{}, this.belief...)

	next.belief[nk] = sampledTruth[nk]
	if next.belief[nk]&Wumpus != 0 {
		t = true
		r = 0
	} else if next.belief[nk]&Pit != 0 {
		t = true
		r = -.01 * float64(1000-this.count)
	} else {
		t = false
		r = -.01
	}

	next.count++

	return
}
func (this Hunter) Shoot() (next Hunter, r float64, t bool) {

	sampledTruth, _ := this.belief.SampleWorld()

	fmt.Println("shooting in\n", sampledTruth)

	xmin, xmax := this.x, this.x
	ymin, ymax := this.y, this.y
	switch this.dir {
	case 0:
		ymin = 0
	case 1:
		xmax = 3
	case 2:
		ymax = 3
	case 3:
		xmin = 0
	}

	r = 0
	for x := xmin; x <= xmax; x++ {
		for y := ymin; y <= ymax; y++ {
			k := x + 4*y
			if sampledTruth[k]&Wumpus != 0 {
				r = 1
			}
		}
	}

	t = true

	return
}

func (this MapBelief) String() (res string) {
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			res += "\t"
			if this[x+4*y] == 0 {
				res += "U"
			}
			if this.GetFlag(x, y, Empty) {
				res += "E"
			}
			if this.GetFlag(x, y, Breezy) {
				res += "B"
			}
			if this.GetFlag(x, y, Stinky) {
				res += "S"
			}
			if this.GetFlag(x, y, Pit) {
				res += "P"
			}
			if this.GetFlag(x, y, Wumpus) {
				res += "W"
			}
		}
		res += "\n"
	}
	return
}

func (this Hunter) String() (res string) {
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			res += "\t"
			if x == int(this.x) && y == int(this.y) {
				switch this.dir {
				case 0:
					res += "^"
				case 1:
					res += ">"
				case 2:
					res += "v"
				case 3:
					res += "<"
				}
			}
			if this.belief[x+4*y] == 0 {
				res += "U"
			}
			if this.belief.GetFlag(x, y, Empty) {
				res += "E"
			}
			if this.belief.GetFlag(x, y, Breezy) {
				res += "B"
			}
			if this.belief.GetFlag(x, y, Stinky) {
				res += "S"
			}
			if this.belief.GetFlag(x, y, Pit) {
				res += "P"
			}
			if this.belief.GetFlag(x, y, Wumpus) {
				res += "W"
			}
		}
		res += "\n"
	}
	return
}
