package wumpus

import (
	"sort"
	"gostat.googlecode.com/hg/stat"
)

const (
	Unknown = 0
	Empty   = 1 << iota
	Breezy
	Stinky
	Pit
	Wumpus

	PPit = .2
)

var SquareMapToValue = sort.IntArray{
	Unknown,
	Empty,
	Breezy,
	Stinky,
	Breezy | Stinky,
	Pit,
	Stinky | Pit,
	Breezy | Stinky | Pit,
	Wumpus,
	Breezy | Wumpus,
	Pit | Wumpus,
}

func GetIndex(status int32) int32 {
	return int32(SquareMapToValue.Search(int(status)))
}
func GetValue(index int32) int32 {
	return int32(SquareMapToValue[index])
}

func flip(p float64) bool {
	return p > stat.NextUniform()
}

type MapBelief []int32

func (this MapBelief) Known(x, y int) bool {
	return this[x+4*y] != 0
}
func (this MapBelief) GetFlag(x, y int, f int32) bool {
	return (this[x+4*y] & f) != 0
}
func (this MapBelief) SetFlag(x, y int, f int32) {
	this[x+4*y] = this[x+4*y] | f
}
func (this MapBelief) SetAdjacent(x, y int, f int32) {
	if x > 0 {
		this.SetFlag(x-1, y, f)
	}
	if x < 3 {
		this.SetFlag(x+1, y, f)
	}
	if y > 0 {
		this.SetFlag(x, y-1, f)
	}
	if y < 3 {
		this.SetFlag(x, y+1, f)
	}
}
func (this MapBelief) AdjacentFlag(x, y int, f int32) (known, set bool) {
	if x > 0 {
		if this.Known(x-1, y) {
			known = true
			if this.GetFlag(x-1, y, f) {
				set = true
			}
		}
	}
	if x < 3 {
		if this.Known(x+1, y) {
			known = true
			if this.GetFlag(x+1, y, f) {
				set = true
			}
		}
	}
	if y > 0 {
		if this.Known(x, y-1) {
			known = true
			if this.GetFlag(x, y-1, f) {
				set = true
			}
		}
	}
	if y < 3 {
		if this.Known(x, y+1) {
			known = true
			if this.GetFlag(x, y+1, f) {
				set = true
			}
		}
	}
	return
}
func (this MapBelief) AdjacentStinky(x, y int) (known bool, count int) {

	if x > 0 {
		if this.Known(x-1, y) {
			known = true
			if this.GetFlag(x-1, y, Stinky) {
				count++
			} else {
				count = 0
				return
			}
		}
	}
	if x < 3 {
		if this.Known(x+1, y) {
			known = true
			if this.GetFlag(x+1, y, Stinky) {
				count++
			} else {
				count = 0
				return
			}
		}
	}
	if y > 0 {
		if this.Known(x, y-1) {
			known = true
			if this.GetFlag(x, y-1, Stinky) {
				count++
			} else {
				count = 0
				return
			}
		}
	}
	if y < 3 {
		if this.Known(x, y+1) {
			known = true
			if this.GetFlag(x, y+1, Stinky) {
				count++
			} else {
				count = 0
				return
			}
		}
	}

	return
}

func (this MapBelief) IsPit(x, y int) bool {
	return this.GetFlag(x, y, Pit)
}

func (this MapBelief) CanBePit(x, y int) (canBePit bool) {
	var known bool
	known, canBePit = this.AdjacentFlag(x, y, Breezy)
	canBePit = !known || canBePit
	return
}
func (this MapBelief) CanBeWumpus(x, y int) (canBeWumpus bool) {
	var known bool
	known, canBeWumpus = this.AdjacentFlag(x, y, Stinky)
	canBeWumpus = !known || canBeWumpus
	return
}

func (this MapBelief) IsPitConsistent(truth MapBelief) bool {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if this.GetFlag(x, y, Breezy) {
				_, pit := truth.AdjacentFlag(x, y, Pit)
				if !pit {
					return false
				}
			}
		}
	}
	return true
}

func (this MapBelief) SampleWorld() (truth MapBelief, rejects int) {
	truth = make(MapBelief, 16)

	var wumpusLocs []int
	smellCount := 0
	anySmell := false
	for k, v := range this {
		if (v & Stinky) != 0 {
			anySmell = true
			smellCount++
		}
		if v&Wumpus != 0 {
			wumpusLocs = []int{k}
		}
	}

	if len(wumpusLocs) == 0 {
		for x := 0; x < 4; x++ {
			for y := 0; y < 4; y++ {
				if x == 0 && y == 0 {
					continue
				}
				known, c := this.AdjacentStinky(x, y)
				if (anySmell && known && c == smellCount) || (!anySmell && !known) {
					wumpusLocs = append(wumpusLocs, x+4*y)
				}
			}
		}
	}
	rejects = 0
	for {
		for i := range truth {
			truth[i] = 0
		}
		for x := 0; x < 4; x++ {
			for y := 0; y < 4; y++ {
				k := x + y*4
				if k == 0 {
					continue
				}
				if this[k] != 0 {
					truth[k] = this[k]
				} else {
					if this.IsPit(x, y) || (this.CanBePit(x, y) && flip(.2)) {
						truth.SetFlag(x, y, Pit)
						truth.SetAdjacent(x, y, Breezy)
					}
				}
			}
		}
		if this.IsPitConsistent(truth) {
			break
		}
		rejects++
	}
	wk := wumpusLocs[stat.NextRange(int64(len(wumpusLocs)))]
	x := wk % 4
	y := wk / 4
	truth.SetFlag(x, y, Wumpus)
	truth.SetAdjacent(x, y, Stinky)

	for k, v := range this {
		if v != 0 {
			truth[k] = v
		}
	}

	for k, v := range truth {
		if v == 0 {
			truth[k] = Empty
		}
	}

	return
}
