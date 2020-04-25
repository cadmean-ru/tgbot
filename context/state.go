package context

//Interface that is used by the bot to preserve state of the user.
//You need to implement it to use scenarios.
type StateProvider interface {
	Save(ctx UpdateContext, state State) error
	Load(ctx UpdateContext) (State, error)
}

//Represents current state of scenario.
type State struct {
	Step     string
	Scenario string
	Data     map[string]interface{}
}

