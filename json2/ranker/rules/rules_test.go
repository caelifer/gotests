package rules_test

import (
	"testing"

	"github.com/caelifer/gotests/json2/ranker/rules"
)

var True = rules.MakeRule(
	"True",
	func(interface{}) bool {
		return true
	},
)

func TestNot(t *testing.T) {
	if res := rules.Not(True).Assert(nil); res != false {
		t.Errorf("[FAILED] got: %v, wanted: %v", res, false)
	}
}

func TestAny(t *testing.T) {
	if res := rules.Any(nil, True, rules.Not(True)); res != true {
		t.Errorf("[FAILED] got: %v, wanted: %v", res, true)
	}
}

func TestAll(t *testing.T) {
	if res := rules.All(nil, True, rules.Not(True)); res != false {
		t.Errorf("[FAILED] got: %v, wanted: %v", res, false)
	}
}
