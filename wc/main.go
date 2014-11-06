package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"test/wc/lex"
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

	fmt.Printf("\t%d\t%d\t%d\t%d\n", lines, words, chars, runes)
	// panic("Testing stack")
	// fmt.Println(str)
}
