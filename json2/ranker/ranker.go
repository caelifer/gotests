package ranker

type Scorer interface {
	Score(interface{}) MatchScore
}

type Ranker struct {
	criteria interface{}
}

func New(c interface{}) *Ranker {
	return &Ranker{c}
}

func (r Ranker) BestMatch(scores ...Scorer) Scorer {
	if len(scores) == 0 {
		return nil
	}

	hs := NoMatch // highest score
	hi := 0       // highest score index

	for i, sc := range scores {
		cs := sc.Score(r.criteria)
		if hs < cs {
			hs = cs
			hi = i
		}
	}

	if hs == NoMatch {
		return nil
	}

	return scores[hi]
}

type MatchScore int

const (
	NoMatch MatchScore = iota - 1
	PartialMatch
	GoodMatch
	PerfectMatch
)
