package coffee

import (
	"gostat.googlecode.com/hg/stat"
	"go-glue.googlecode.com/hg/rlglue"
	"go-glue.googlecode.com/hg/rltools/discrete"
)

type Config struct {
	RewardDry, RewardUserHasCoffee																float64
	ProbGetWetUmbrella, ProbGetWetNoUmbrella, ProbMove, ProbDeliverCoffee, ProbKeepCoffeeNoOffice, ProbKeepCoffeeOffice, ProbFindUmbrella, ProbBuyCoffee	float64
	ProbDrinkCoffee																		float64
	DiscountFactor																		float64
}

func ConfigDefault() (cfg Config) {
	cfg.RewardDry = .1
	cfg.RewardUserHasCoffee = .9
	cfg.ProbGetWetNoUmbrella = .9
	cfg.ProbGetWetUmbrella = .1
	cfg.ProbMove = .9
	cfg.ProbDeliverCoffee = .8
	cfg.ProbKeepCoffeeNoOffice = .2
	cfg.ProbKeepCoffeeOffice = .1
	cfg.ProbFindUmbrella = .9
	cfg.ProbBuyCoffee = .9
	cfg.DiscountFactor = 0.9
	return
}

const (
	Move	= iota
	BuyCoffee
	GetUmbrella
	DeliverCoffee
)

type Env struct {
	*discrete.MDPEnvironment
	mdp	*MDP
}

func NewEnv(cfg Config) (this *Env) {
	this = new(Env)
	this.mdp = NewMDP(cfg)
	this.MDPEnvironment = discrete.NewMDPEnvironment(this.mdp, this.mdp.GetTask(), 0)
	return
}
func (this *Env) EnvStart() (obs rlglue.Observation) {
	startState := uint64(stat.NextRange(int64(this.mdp.GetTask().Obs.Ints.Count())))
	obs = rlglue.NewObservation(this.mdp.GetTask().Obs.Ints.Values(startState), []float64{}, []byte{})
	this.LastState = startState
	return
}
