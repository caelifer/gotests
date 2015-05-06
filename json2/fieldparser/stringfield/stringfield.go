package stringfield

import (
	"bytes"
	"encoding/json"

	"github.com/caelifer/gotests/json2/fieldparser"
	"github.com/caelifer/gotests/json2/ranker"
	"github.com/caelifer/gotests/json2/ranker/rules"
)

type StringFieldParser struct {
	name  string
	rules []rules.Rule
}

func (sp StringFieldParser) Score(cond interface{}) ranker.MatchScore {
	if rules.All(cond, sp.rules...) {
		return ranker.PartialMatch
	}

	return ranker.NoMatch
}

func (sp StringFieldParser) Parse(meta fieldparser.Meta, jsn []byte) (string, error) {
	var key, value string

	// Split on key/value
	rdata := bytes.SplitN(jsn, []byte(":"), 2)

	// Parse JSON chunk

	// Key
	if err := json.Unmarshal(rdata[0], &key); err != nil {
		return "", err
	}

	// Value
	if err := json.Unmarshal(rdata[1], &value); err != nil {
		return "", err
	}

	// Make sure key is matching
	if key == meta.ID() {
		return value, nil
	}

	// Bad meta
	return "", fieldparser.BadMetaError
}

func init() {
	fieldparser.RegisterParser(
		StringFieldParser{
			name: "StringFieldParser",
			// Add rules
			rules: append([]rules.Rule{},
				rules.MakeRule(
					"StringFieldType",
					func(cond interface{}) bool {
						// Make sure to handle run-time type narrowing properly
						if meta, ok := cond.(fieldparser.Meta); ok {
							return meta.Type() == "string"
						}
						return false // Failed meta test
					},
				),
			),
		},
	)
}
