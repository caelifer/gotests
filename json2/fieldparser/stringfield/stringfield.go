package stringfield

import (
	"github.com/caelifer/gotests/json2/fieldparser"
	"github.com/caelifer/gotests/json2/ranker"
	"github.com/caelifer/gotests/json2/ranker/rules"
)

type StringFieldParser struct {
	name  string
	rules []rules.Rule
}

func (sp StringFieldParser) Score(cond interface{}) ranker.MatchScore {
	if rules.All(cond.(fieldparser.Meta), sp.rules...) { // assert correct type
		return ranker.PartialMatch
	}

	return ranker.NoMatch
}

func (sp StringFieldParser) Parse(jsn []byte) (string, error) {
	return "DUMMY", nil
}

func init() {
	fieldparser.RegisterParser(StringFieldParser{
		name: "StringFieldParser",
		// Add rules
		rules: append([]rules.Rule{},
			rules.MakeRule(
				"StringFieldType",
				func(cond interface{}) bool {
					return cond.(fieldparser.Meta).Type() == "string"
				},
			),
		),
	})
}
