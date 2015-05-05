package main

import "fmt"

// package rank

type Criteria interface{}

type Scorer interface {
	Score(Criteria) int
}

type Ranker struct {
	criteria Criteria
}

func (r Ranker) BestMatch(scores ...Scorer) Scorer {
	if len(scores) == 0 {
		return nil
	}

	hs := 0 // heiest score
	hi := 0 // heiest score index

	for i, sc := range scores {
		cs := sc.Score(r.criteria)
		if hs < cs {
			hs = cs
			hi = i
		}
	}
	return scores[hi]

}

type score int

func (s score) Score(Criteria) int {
	return int(s)
}

func main() {
	fmt.Println(new(Ranker).BestMatch(
		score(1),
		score(2),
		score(3),
		score(4),
		score(5),
		score(5),
		score(6),
	))

	rule1 := MakeRule("Name", func(meta interface{}) bool {
		return "test1" == meta.(Meta).Name
	})

	rule2 := MakeRule("Type", func(meta interface{}) bool {
		return "string" == meta.(Meta).Type
	})

	rule3 := MakeRule("Type", func(meta interface{}) bool {
		return "string2" == meta.(Meta).Type
	})

	m := Meta{
		Name: "test1",
		Type: "string",
	}

	fmt.Println(
		m.All(rule1, rule3),
		m.Any(rule1, rule3),
		m.Not(rule2),
	)
}

type Rule struct {
	Name string
	Eval func(interface{}) bool
}

func MakeRule(name string, test func(interface{}) bool) *Rule {
	r := new(Rule) // alloc on heap

	r.Name = name
	r.Eval = func(data interface{}) bool {
		return test(data)
	}

	return r
}

type Meta struct {
	Name, Type string
}

func (meta Meta) All(rules ...*Rule) bool {
	res := false
	for _, r := range rules {
		if res = r.Eval(meta); !res {
			return false
		}
	}
	return res
}

func (meta Meta) Any(rules ...*Rule) bool {

	for _, r := range rules {
		if r.Eval(meta) {
			return true
		}
	}
	return false
}

func (meta Meta) Not(rule *Rule) bool {
	return !rule.Eval(meta)
}
