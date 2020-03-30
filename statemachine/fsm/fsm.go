package fsm

import (
	"fmt"
	"sync"

	"github.com/caelifer/gotests/statemachine/event"
	"github.com/caelifer/gotests/statemachine/state"
)

// FSM is an interface that represents Finite State Machine operations
type FSM interface {
	Process(event.Event) error
	State() state.State
}

func New(s state.State) FSM {
	return &fsm{current: s}
}

type fsm struct {
	sync.RWMutex
	current state.State
}

func (sm *fsm) Process(evt event.Event) error {
	newState, err := sm.current.TransitionTo(evt.TargetState())
	if err != nil {
		return fmt.Errorf("bad state transition: %v -> %v: %v", sm.current, newState, err)
	}

	sm.Lock() // Take exclusive lock
	sm.current = newState
	sm.Unlock()
	return nil
}

func (sm *fsm) State() state.State {
	sm.RLock()
	defer sm.RUnlock()
	return sm.current
}
