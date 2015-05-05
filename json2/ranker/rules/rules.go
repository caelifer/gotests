package rules

type Rule interface {
	Name() string
	Assert(interface{}) bool
}

func Not(r Rule) Rule {
	assert := func(cond interface{}) bool {
		return !r.Assert(cond)
	}
	return MakeRule(r.Name(), assert)
}

func All(cond interface{}, rules ...Rule) bool {
	res := false
	for _, r := range rules {
		if res = r.Assert(cond); !res {
			return res
		}
	}
	return res
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
