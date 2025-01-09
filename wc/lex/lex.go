package lex

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	Error TokenType = iota
	Character
	Line
	Word
	EOF
)

func (tt TokenType) String() string {
	res := "NOOP"

	switch tt {
	case Character:
		res = "Character"
	case Line:
		res = "Line"
	case Word:
		res = "Word"
	case EOF:
		res = "<EOF>"
	default:
		res = "Error"
	}

	return res
}

type Token struct {
	Type  TokenType
	Value string
	// Len   int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %q, Value: %q}", t.Type, t.Value)
}

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input             string
	start, pos, width int
	tokens            chan Token
	initFn            stateFn
}

func Lex(input string) <-chan Token {
	out := make(chan Token)
	lex := &Lexer{
		input:  input,
		tokens: out,
		initFn: lexText,
	}

	// Start lexing on the background
	go lex.run()

	return out
}

// FSA state machine
func (l *Lexer) run() {
	for s := lexText; s != nil; {
		s = s(l)
	}
}

// lexText is a basic scanner for Characters and Line. It will also transition to
// EOF and Word scanner
func lexText(l *Lexer) stateFn {
	// log.Println("Entering lexText")
	sym := l.next()
	if sym == eof {
		l.emit(EOF)
		return nil
	}

	// See if this is a start of word
	if unicode.IsDigit(sym) || unicode.IsLetter(sym) || sym == '_' {
		l.putback()
		return lexWord
	}

	// Check for the EOL
	if sym == '\n' {
		l.emit(Line)
	}

	l.emit(Character)
	return lexText
}

// lexWord is a word scanner. Word is considered to be a sequence of runes
// starting with Number, Letter or Underscore and followed by series Letter,
// Number, Underscore or '-'
func lexWord(l *Lexer) stateFn {
	// log.Println("Entering lexWord")
	// We already know that the first character is good
	sym := l.next()
	_ = sym // quelle linter warning
	l.emit(Character)

	for {
		sym = l.next()
		if sym == eof {
			if l.start < l.pos {
				l.emit(Word)
			}
			return lexEOF
		}

		if !unicode.IsDigit(sym) && !unicode.IsLetter(sym) && sym != '-' && sym != '—' && sym != '_' {
			l.putback()
			break
		}

		l.emit(Character)
	}

	l.emit(Word)

	return lexText
}

// Properly handle EOF
func lexEOF(l *Lexer) stateFn {
	l.emit(EOF)
	close(l.tokens)
	return nil
}

func (l *Lexer) emit(tt TokenType) {
	tok := Token{
		Type: tt,
	}

	if l.pos > l.start {
		tok.Value = l.input[l.start:l.pos]
		// tok.Len = len(l.input[l.start:l.pos])
		l.start = l.pos
	}

	// log.Printf("In emit with %v", tok)

	l.tokens <- tok
}

const eof = rune(0)

func (l *Lexer) next() rune {
	r := eof
	if l.start < len(l.input) {
		r, l.width = utf8.DecodeRuneInString(l.input[l.start:])
		// log.Printf("next() got %q", string(r))
		l.pos += l.width
	}
	return r
}

func (l *Lexer) putback() {
	l.pos -= l.width
	l.width = 0
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.putback()
	return r
}

// vim: :ts=4:sw=4:noexpandtab:ai
