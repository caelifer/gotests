package strings

import "testing"

func Test(t *testing.T) {
	ttable := map[string]string{
		"":                "",
		"a":               "a",
		"aa":              "aa",
		"ab":              "ba",
		"Hello, gopher!":  "!rehpog ,olleH",
		"Здравствуй, ГО!": "!ОГ ,йувтсвардЗ",
	}

	for test, want := range ttable {
		got := Reverse(test)
		if got != want {
			t.Errorf("For '%q' got '%q', expected '%q'!", test, got, want)
		}
	}
}
