package fieldparser_test

import (
	"testing"

	"github.com/caelifer/gotests/json2/fieldparser"

	// Load all available field type parsers
	_ "github.com/caelifer/gotests/json2/fieldparser/numberfield"
	_ "github.com/caelifer/gotests/json2/fieldparser/stringfield"
)

func TestFieldParsers(t *testing.T) {
	tests := []struct {
		meta fieldparser.Meta
		jsn  []byte
		want string
		err  error
	}{
		{
			meta: meta{"dummy", "number"},
			jsn:  []byte(`"dummy": 1234.5`),
			want: "1234.5",
			err:  nil,
		},
		{
			meta: meta{"dummy", "string"},
			jsn:  []byte(`"dummy": "dummytest"`),
			want: "dummytest",
			err:  nil,
		},
	}

	for _, tst := range tests {
		got, err := fieldparser.ParseJiraField(tst.meta, tst.jsn)

		// Check error first
		if err != nil {
			if tst.err == nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}

		// Check return value
		if tst.want != got {
			t.Errorf("Got: %q, wanted: %q for input: %s", got, tst.want, string(tst.jsn))
		}
	}
}

// Mock type that implements fieldparser.Meta interface
type meta struct {
	id, kind string
}

func (m meta) ID() string {
	return m.id
}

func (m meta) Type() string {
	return m.kind
}

func (m meta) Items() string {
	return ""
}

func (m meta) Custom() string {
	return ""
}
