package numberfield

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/caelifer/gotests/json2/fieldparser"
	"github.com/caelifer/gotests/json2/ranker"
	"github.com/caelifer/gotests/json2/ranker/rules"
)

type NumberFieldParser struct {
	name  string
	rules []rules.Rule
}

func (sp NumberFieldParser) Score(cond interface{}) ranker.MatchScore {
	if rules.All(cond, sp.rules...) {
		return ranker.PartialMatch
	}

	return ranker.NoMatch
}

func (sp NumberFieldParser) Parse(meta fieldparser.Meta, data []byte) (string, error) {
	var (
		key   string
		value float64
	)

	// Split on key/value
	parts := bytes.SplitN(data, []byte(":"), 2)

	// Parse JSON chunk

	// Key
	if err := json.Unmarshal(parts[0], &key); err != nil {
		return "", err
	}

	// Value
	if err := json.Unmarshal(parts[1], &value); err != nil {
		return "", err
	}

	// Make sure key is matching
	if key == meta.ID() {
		return fmt.Sprintf("%g", value), nil
	}

	// Bad meta
	return "", fieldparser.BadMetaError
}

func init() {
	fieldparser.RegisterParser(
		NumberFieldParser{
			name: "NumberFieldParser",
			// Add rules
			rules: append([]rules.Rule{},
				rules.MakeRule(
					"NumberFieldType",
					func(cond interface{}) bool {
						// Make sure to handle run-time type narrowing properly
						if meta, ok := cond.(fieldparser.Meta); ok {
							return meta.Type() == "number"
						}
						return false // Failed meta test
					},
				),
			),
		},
	)
}
