package numberfield_test

import (
	"testing"

	"github.com/caelifer/gotests/json2/fieldparser"

	// Load field type parser
	_ "github.com/caelifer/gotests/json2/fieldparser/numberfield"
)

func TestNumberFieldParserError(t *testing.T) {
	tests := []struct {
		meta fieldparser.Meta
		jsn  []byte
		want string
		err  error
	}{
		{
			meta: meta{"dummy", "number2"},
			jsn:  []byte(`"dummy": 1234.5`),
			want: "",
			err:  fieldparser.NoMatchingParserError,
		},
	}

	for _, tst := range tests {
		_, err := fieldparser.ParseJiraField(tst.meta, tst.jsn)

		// Check error first
		if err != nil {
			if tst.err != err {
				t.Errorf("Unexpected error: %v", err)
			}
		} else {
			t.Errorf("Unexpected success -  got %v, wanted: %v for input: %s", err, tst.err, string(tst.jsn))
		}
	}
}

func TestNumberFieldParser(t *testing.T) {
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

// Test structure that implements fieldparser.Meta interface
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
