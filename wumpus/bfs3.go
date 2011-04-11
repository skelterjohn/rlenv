package wumpus

import (
	"github.com/skelterjohn/rlalg/bfs3"
	"github.com/skelterjohn/rlbayes"
	"go-glue.googlecode.com/hg/rlglue"
)

func GetConfiguredPrior() bfs3.Prior {
	return func(task *rlglue.TaskSpec) (prior bayes.BeliefState) {
		prior = NewBelief(make(MapBelief, 16))
		return
	}
}
func NewBFS3Agent(bfs3cfg bfs3.Config) (agent *bfs3.BFS3Agent) {
	prior := GetConfiguredPrior()
	agent = bfs3.New(prior)
	agent.Cfg = bfs3cfg
	return
}
