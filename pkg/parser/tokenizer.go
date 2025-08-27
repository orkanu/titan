package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type tokenType int

const (
	tEOF tokenType = iota
	tIdent
	tString
	tNumber
	tOp
	tAnd
	tOr
	tLParen
	tRParen
)

type token struct {
	typ tokenType
	val string
}

type lexer struct {
	input []rune
	pos   int
}

func newLexer(s string) *lexer {
	return &lexer{input: []rune(s)}
}

func (l *lexer) next() token {
	for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.pos++
	}
	if l.pos >= len(l.input) {
		return token{typ: tEOF}
	}

	ch := l.input[l.pos]

	// Ident or variable
	if unicode.IsLetter(ch) || ch == '_' {
		start := l.pos
		for l.pos < len(l.input) && (unicode.IsLetter(l.input[l.pos]) || unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '_') {
			l.pos++
		}
		word := string(l.input[start:l.pos])
		return token{typ: tIdent, val: word}
	}

	// String literal
	if ch == '\'' || ch == '"' {
		quote := ch
		l.pos++
		start := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != quote {
			l.pos++
		}
		val := string(l.input[start:l.pos])
		l.pos++ // consume closing
		return token{typ: tString, val: val}
	}

	// Number
	if unicode.IsDigit(ch) {
		start := l.pos
		for l.pos < len(l.input) && (unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '.') {
			l.pos++
		}
		return token{typ: tNumber, val: string(l.input[start:l.pos])}
	}

	// Operators
	switch {
	case strings.HasPrefix(string(l.input[l.pos:]), "&&"):
		l.pos += 2
		return token{typ: tAnd, val: "&&"}
	case strings.HasPrefix(string(l.input[l.pos:]), "||"):
		l.pos += 2
		return token{typ: tOr, val: "||"}
	case strings.HasPrefix(string(l.input[l.pos:]), "=="):
		l.pos += 2
		return token{typ: tOp, val: "=="}
	case strings.HasPrefix(string(l.input[l.pos:]), "!="):
		l.pos += 2
		return token{typ: tOp, val: "!="}
	case strings.HasPrefix(string(l.input[l.pos:]), ">="):
		l.pos += 2
		return token{typ: tOp, val: ">="}
	case strings.HasPrefix(string(l.input[l.pos:]), "<="):
		l.pos += 2
		return token{typ: tOp, val: "<="}
	case ch == '>':
		l.pos++
		return token{typ: tOp, val: ">"}
	case ch == '<':
		l.pos++
		return token{typ: tOp, val: "<"}
	case ch == '(':
		l.pos++
		return token{typ: tLParen, val: "("}
	case ch == ')':
		l.pos++
		return token{typ: tRParen, val: ")"}
	default:
		panic(fmt.Sprintf("unexpected character: %q", ch))
	}
}
