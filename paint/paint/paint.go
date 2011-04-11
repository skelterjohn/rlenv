package main

import (
	"fmt"
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"go-glue.googlecode.com/hg/rlglue"
	"github.com/skelterjohn/rlenv/paint"
)

func main() {
	defer nicetrace.Print()
	config := paint.ConfigDefault()
	argcfg.LoadArgs(&config)
	env := paint.New(config)
	env.Live = false
	if err := rlglue.LoadEnvironment(env); err != nil {
		fmt.Println("Error running paint: %v\n", err)
	}
}
