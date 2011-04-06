package wumpus

import (
	"go-glue.googlecode.com/hg/rlglue"
)

func Main() {
	env := New()
	if err := rlglue.LoadEnvironment(env); err != nil {
		println(err.String())
	}
}
