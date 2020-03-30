package event

import "github.com/caelifer/gotests/statemachine/state"

type Event interface {
	TargetState() state.State
}
