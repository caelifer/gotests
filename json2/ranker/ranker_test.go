package ranker_test

import (
	"testing"

	"github.com/caelifer/gotests/json2/ranker"
)

func TestNewRanker(t *testing.T) {
	_ = ranker.New(nil)
}

func TestBestMatch(t *testing.T) {
	tests := []struct {
		scores []ranker.Scorer
		want   ranker.Scorer
	}{
		{
			scores: []ranker.Scorer{},
			want:   nil,
		},
		{
			scores: []ranker.Scorer{score(ranker.NoMatch)},
			want:   nil,
		},
		{
			scores: []ranker.Scorer{score(ranker.PartialMatch)},
			want:   score(ranker.PartialMatch),
		},
		{
			scores: []ranker.Scorer{score(ranker.GoodMatch), score(ranker.PartialMatch)},
			want:   score(ranker.GoodMatch),
		},
		{
			scores: []ranker.Scorer{score(ranker.GoodMatch), score(ranker.PartialMatch), score(ranker.PerfectMatch)},
			want:   score(ranker.PerfectMatch),
		},
	}

	for _, tst := range tests {
		if got := ranker.New(nil).BestMatch(tst.scores...); got != tst.want {
			t.Errorf("[FAILED] BestMatch - wanted: %d, got: %d on input: %+v", got, tst.want, tst.scores)
		}
	}
}

// test implementation of ranker.Scorer
type score ranker.MatchScore

func (s score) Score(interface{}) ranker.MatchScore {
	return ranker.MatchScore(s)
}
