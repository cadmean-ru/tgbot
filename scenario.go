package tgbot


import (
	"github.com/cadmean-ru/tgbot/context"
)

//The first return should be the name of the next step or "" for the scenario to be terminated.
//The error is not nil, error handler will be called.
type StepHandler func(ctx *context.UpdateContext) (string, error)

//Scenario (or use case, or user flow)
//Consists of several steps. A step consists of user input and bot response.
//The first step is triggered when the user sends on of the trigger commands or phrases.
//Each step has a handler, that should return the name of the next step to be executed after the next user input.
type Scenario struct {
	Triggers []string
	Steps    []Step
	Name     string
}

type Step struct {
	Name    string
	Handler StepHandler
}

//Used to conveniently build new scenarios.
type Builder struct {
	scenario *Scenario
}

//Start building new scenario with name
//The name of scenario should be uniq
func NewScenario(name string) *Builder  {
	return &Builder{
		scenario: &Scenario{
			Name:     name,
			Steps:    make([]Step, 0),
		},
	}
}

//Sets commands or messages when the first step of the scenario is triggered
func (b *Builder) TriggeredBy(triggers ...string) *Builder {
	b.scenario.Triggers = triggers
	return b
}

//Adds a new step to the scenario.
//Step name should be uniq.
func (b *Builder) AddStep(name string, handler StepHandler) *Builder {
	b.scenario.Steps = append(b.scenario.Steps, Step{
		Name:    name,
		Handler: handler,
	})
	return b
}

//Returns the newly constructed scenario
func (b *Builder) Create() *Scenario {
	return b.scenario
}