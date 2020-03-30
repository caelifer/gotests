package state

type State interface {
	TransitionTo(State) (State, error)
}
