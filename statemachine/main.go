package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/caelifer/gotests/statemachine/fsm"
	"github.com/caelifer/gotests/statemachine/state"
)

type OrderState int

const Init OrderState = 0
const (
	New OrderState = 1 << iota
	Acked
	PartiallyFilled
	Filled
	Rejected
	Canceled
	Fulfilled
)

type ostm int // Order State Transition Map

var stateMapRules = []struct {
	s  OrderState
	sm ostm
}{
	{Init, ostm(New)},
	{New, ostm(Acked | Rejected)},
	{Acked, ostm(PartiallyFilled | Filled | Rejected | Canceled)},
	{PartiallyFilled, ostm(PartiallyFilled | Filled | Rejected | Canceled)},
	{Filled, ostm(Fulfilled)},
	{Canceled, ostm(Fulfilled)},
	{Rejected, ostm(Fulfilled)},
}

func lookup(os OrderState) ostm {
	for _, smr := range stateMapRules {
		if smr.s == os {
			return smr.sm
		}
	}
	return ostm(0)
}

func (tm ostm) isValid(os OrderState) bool {
	return tm&ostm(os) != 0
}

// Implement state.State interface

func (os OrderState) TransitionTo(newState state.State) (state.State, error) {
	sm := lookup(os)
	if sm.isValid(newState.(OrderState)) {
		fmt.Printf("transitioning %v -> %v\n", os, newState)
		return newState, nil
	}

	return nil, fmt.Errorf("transition not allowed")
}

func (os OrderState) String() string {
	switch os {
	case Init:
		return "init"
	case New:
		return "new"
	case Acked:
		return "acked"
	case PartiallyFilled:
		return "pfilled"
	case Filled:
		return "filled"
	case Rejected:
		return "rejected"
	case Canceled:
		return "canceled"
	case Fulfilled:
		return "fulfilled"
	}
	return "invalid"
}

// Implement event.Event interface
type OrderMessage struct {
	id          int
	name        string
	payload     []byte
	targetState OrderState
}

func NewOrderMessage(id int, name string, targetState OrderState) OrderMessage {
	return OrderMessage{
		id:          id,
		name:        name,
		targetState: targetState,
	}
}

func (om OrderMessage) TargetState() state.State {
	return om.targetState
}

func (om OrderMessage) String() string {
	return fmt.Sprintf("OrderMessage{id: %d, type: %q, targetState: %v}", om.id, om.name, om.targetState)
}

// State function
type stateFn func(OrderMessage) stateFn

///////////////// Main program entry ////////////////

func main() {
	rand.Seed(time.Now().Unix())

	id := 1

	// Order FSM
	osm := fsm.New(Init)

	n := 0
LOOP:
	for {
		var om OrderMessage

		switch osm.State() {
		case Init:
			om = NewOrderMessage(id, "NewOrderSingle", New)
		case New:
			om = NewOrderMessage(id, "OrderAck", Acked)
		case Acked:
			choices := []OrderState{PartiallyFilled, PartiallyFilled, PartiallyFilled, PartiallyFilled, Filled, Init, Canceled, Rejected}
			// ns := choices[rand.Int()%len(choices)]
			ns := choices[n]
			n++
			switch ns {
			case Canceled:
				om = NewOrderMessage(id, "OrderCancel", ns)
			case Rejected:
				om = NewOrderMessage(id, "OrderReject", ns)
			case PartiallyFilled:
				om = NewOrderMessage(id, "OrderPartialFill", ns)
			case Filled:
				om = NewOrderMessage(id, "OrderFill", ns)
			case Init:
				om = NewOrderMessage(id, "BadMessage", ns)
			}
		case PartiallyFilled:
			choices := []OrderState{PartiallyFilled, PartiallyFilled, PartiallyFilled, PartiallyFilled, Filled, Init, Canceled, Rejected}
			ns := choices[n]
			n++
			switch ns {
			case Canceled:
				om = NewOrderMessage(id, "OrderCancel", ns)
			case Rejected:
				om = NewOrderMessage(id, "OrderReject", ns)
			case PartiallyFilled:
				om = NewOrderMessage(id, "OrderPartialFill", ns)
			case Filled:
				om = NewOrderMessage(id, "OrderFill", ns)
			case Init:
				om = NewOrderMessage(id, "BadMessage", ns)
			}
		case Canceled, Rejected, Filled:
			om = NewOrderMessage(id, "OrderFulfilled", Fulfilled)
		case Fulfilled:
			fmt.Printf("OMS: reached the the leaf state: %v\n", osm.State())
			break LOOP
			// default:
			// 	log.Fatalf("Unknown state: %v: %v", osm.State(), om)
		}
		fmt.Println("OMS: received event:", om)
		if err := osm.Process(om); err != nil {
			log.Fatal(err)
		}
	}
}
