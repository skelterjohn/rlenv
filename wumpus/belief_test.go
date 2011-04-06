package wumpus

import (
	"testing"
	"fmt"
	"go-glue.googlecode.com/hg/rlglue"
	"gonicetrace.googlecode.com/hg/nicetrace"
	"gostat.googlecode.com/hg/stat"
)

func TestFunny(t *testing.T) {
	if true {
		return
	}
	stat.TimeSeed()
	belief := MapBelief{Empty, Empty, Unknown, Unknown, Stinky, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown}
	const SS = 10
	c := 0.0
	for i := 0; i < SS; i++ {
		b, _ := belief.SampleWorld()
		if b.GetFlag(0, 2, Wumpus) {
			c++
		}
		fmt.Println(b)
	}
	fmt.Println(c / SS)
}
func TestSample(t *testing.T) {
	if true {
		return
	}
	stat.TimeSeed()
	belief := MapBelief{Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown, Unknown}
	tr := 0
	for i := 0; i < 1000; i++ {
		b, _ := belief.SampleWorld()
		for k := range b {
			if b[k]&Wumpus != 0 || flip(.5) {
				b[k] = 0
			}
		}
		_, r := b.SampleWorld()
		tr += r
	}
	fmt.Printf("%d avg rejects\n", tr/1000)
}
func TestStateCount(t *testing.T) {
}
func TestEnvTraj(t *testing.T) {
	if false {
		return
	}
	stat.TimeSeed()
	defer nicetrace.Print()
	env := New()
	obsi := env.EnvStart()
	indexi := task.Obs.Ints.Index(obsi.Ints())
	b := NewBelief(make(MapBelief, 16))
	b.Teleport(indexi)
	fmt.Printf("%v\n", b.Hunter)
	do := func(what int32) bool {
		action := rlglue.NewAction([]int32{what}, []float64{}, []byte{})
		obs, r, t := env.EnvStep(action)
		fmt.Println(what, r)
		if t {
			fmt.Println("done")
			return false
		}
		index := task.Obs.Ints.Index(obs.Ints())
		bs := b.Update(uint64(what), index, r)
		b = bs.(*Belief)
		fmt.Printf("%v\n", b.Hunter)
		return true
	}
	guess := func() (what int32) {
		fmt.Scanf("%d", &what)
		return
	}
	for do(guess()) {
	}
}
func TestTraj(t *testing.T) {
	if true {
		return
	}
	observed := make(MapBelief, 16)
	truth, _ := observed.SampleWorld()
	observed[0] = truth[0]
	b := NewBelief(observed)
	fmt.Println(b.Hunter)
	do := func(what uint64) bool {
		o, r := b.Next(what)
		fmt.Println(what, r)
		if o == nil {
			fmt.Println("done")
			return false
		}
		b = o.(*Belief)
		fmt.Println(b.Hunter)
		return true
	}
	guess := func() (what uint64) {
		fmt.Scanf("%d", &what)
		return
	}
	for do(guess()) {
	}
}
