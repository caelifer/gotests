package rules

type Rule interface {
	Name() string
	Assert(interface{}) bool
}

func Not(r Rule) Rule {
	return MakeRule(
		r.Name(),
		func(cond interface{}) bool {
			return !r.Assert(cond)
		},
	)
}

func All(cond interface{}, rules ...Rule) bool {
	res := false
	for _, r := range rules {
		// log.Printf("Processing rule: %v", r.Name())
		if res = r.Assert(cond); !res {
			return res
		}
	}
	return res
}

func Any(cond interface{}, rules ...Rule) bool {
	for _, r := range rules {
		if r.Assert(cond) {
			return true
		}
	}
	return false
}

func MakeRule(name string, assert func(interface{}) bool) Rule {
	return rule{
		name:   name,
		assert: assert,
	}
}

// Implementation
type rule struct {
	name   string
	assert func(interface{}) bool
}

func (r rule) Name() string {
	return r.name
}

func (r rule) Assert(cond interface{}) bool {
	return r.assert(cond)
}
