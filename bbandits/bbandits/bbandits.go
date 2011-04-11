package main

import (
	"fmt"
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"go-glue.googlecode.com/hg/rlglue"
	"github.com/skelterjohn/rlenv/bbandits"
)

type Config struct{ NumActions int }

func main() {
	defer nicetrace.Print()
	var cfg Config
	cfg.NumActions = 5
	argcfg.LoadArgs(&cfg)
	ones := make([]float64, cfg.NumActions)
	for i, _ := range ones {
		ones[i] = 1
	}
	genv := bbandits.NewEnv(bbandits.NewBelief(ones, ones))
	if err := rlglue.LoadEnvironment(genv); err != nil {
		fmt.Println("Error running grid: %v\n", err)
	}
}
