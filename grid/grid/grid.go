package main

import (
	"fmt"
	"gonicetrace.googlecode.com/hg/nicetrace"
	"goargcfg.googlecode.com/hg/argcfg"
	"go-glue.googlecode.com/hg/rlglue"
	"github.com/skelterjohn/rlenv/grid"
)

type Config struct{ Width, Height int }

func main() {
	defer nicetrace.Print()
	var config Config
	config.Width, config.Height = 3, 3
	argcfg.LoadArgs(&config)
	genv := grid.New(int32(config.Width), int32(config.Height))
	if err := rlglue.LoadEnvironment(genv); err != nil {
		fmt.Println("Error running grid: %v\n", err)
	}
}
