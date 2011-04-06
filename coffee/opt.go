package coffee

import (
	"go-glue.googlecode.com/hg/rlglue"
	"go-glue.googlecode.com/hg/rlalg/vi"
	"go-glue.googlecode.com/hg/rltools/discrete"
)

type OptAgent struct {
	Cfg	Config
	task	*rlglue.TaskSpec
	mdp	*MDP
	qt	*discrete.QTable
}

func NewOptAgent(Cfg Config) (this *OptAgent) {
	this = new(OptAgent)
	this.Cfg = Cfg
	return
}
func (this *OptAgent) AgentInit(taskString string) {
	this.task, _ = rlglue.ParseTaskSpec(taskString)
	this.mdp = NewMDP(this.Cfg)
	this.qt = discrete.NewQTable(this.task.Obs.Ints.Count(), this.task.Act.Ints.Count())
	vi.ValueIteration(this.qt, this.mdp, 0.01)
}
func (this *OptAgent) AgentStart(obs rlglue.Observation) (act rlglue.Action) {
	s := this.task.Obs.Ints.Index(obs.Ints())
	a := this.qt.Pi(s)
	act = rlglue.NewAction([]int32{int32(a)}, []float64{}, []byte{})
	return
}
func (this *OptAgent) AgentStep(reward float64, obs rlglue.Observation) (act rlglue.Action) {
	s := this.task.Obs.Ints.Index(obs.Ints())
	a := this.qt.Pi(s)
	act = rlglue.NewAction([]int32{int32(a)}, []float64{}, []byte{})
	return
}
func (this *OptAgent) AgentEnd(reward float64) {
	return
}
func (this *OptAgent) AgentCleanup() {
	return
}
func (this *OptAgent) AgentMessage(message string) (reply string) {
	return
}
