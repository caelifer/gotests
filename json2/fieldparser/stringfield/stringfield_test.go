package stringfield_test

import (
	"testing"

	"github.com/caelifer/gotests/json2/fieldparser"

	// Load required field parser
	_ "github.com/caelifer/gotests/json2/fieldparser/stringfield"
)

func TestStringFieldParserError(t *testing.T) {
	tests := []struct {
		meta fieldparser.Meta
		jsn  []byte
		want string
		err  error
	}{
		{
			meta: meta{"dummy", "string2"},
			jsn:  []byte(`"dummy": "dummytest"`),
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
			t.Errorf("Unexpected success -  got %v, wanted: %q for input: %s", err, tst.err, string(tst.jsn))
		}
	}
}

func TestStringFieldParser(t *testing.T) {
	tests := []struct {
		meta fieldparser.Meta
		jsn  []byte
		want string
		err  error
	}{
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
			} else {
				t.Logf("Captured error: %v", err)
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
