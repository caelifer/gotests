package main

import (
	"testing"

	"github.com/caelifer/gotests/json/parser"
)

var tests = []struct {
	test, meta, want string
}{
	{
		test: `"string1": "string1val"`,
		meta: "stringField",
		want: "string1val",
	},
	{
		test: `
	"arrofstrings": [
		"aos1",
		"aos2",
		"aos3"
	]`,
		meta: "arrayOfStrings",
		want: "aos1,aos2,aos3",
	},
	{
		test: `
	"arrofobjecs": [
		{
			"of1": "of1val1",
			"of2": "of2val1",
			"of3": "of3val1"
		},
		{
			"of1": "of1val2",
			"of2": "of2val2"
		}
	]
`,
		meta: "arrayOfObjects",
		want: "{of1:of1val1,of2:of2val1,of3:of3val1},{of1:of1val2,of2:of2val2}",
	},
}

func TestCustJSONParsing(t *testing.T) {
	// Run out test cases
	for _, tst := range tests {
		got := parser.Parse(tst.meta, tst.test)

		if tst.want != got {
			t.Errorf("Got: %q, expected: %q\n", got, tst.want)
		} else {
			t.Logf("[PASS] got %q while parsing %s", got, tst.meta)
		}
	}
}
