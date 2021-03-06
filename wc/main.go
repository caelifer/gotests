package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/caelifer/gotests/wc/lex"
)

func main() {
	var lines, words, runes, chars int
	// str := ""

	content, _ := ioutil.ReadAll(os.Stdin)

	in := lex.Lex(string(content))

LOOP:
	for t := range in {
		switch t.Type {
		case lex.Line:
			lines++
			chars++
			runes++
			// str += t.Value
		case lex.Word:
			words++
		case lex.Character:
			runes++
			chars += len(t.Value)
			// str += t.Value
		default:
			break LOOP
		}
	}

	fmt.Printf("%8d %7d %7d %7d\n", lines, words, chars, runes)
	// panic("Testing stack")
	// fmt.Println(str)
}
