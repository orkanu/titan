package parser

import (
	"fmt"
	"strconv"
)

type parser struct {
	lex  *lexer
	peek token
	ctx  map[string]any
}

func NewParser(expr string, ctx map[string]any) *parser {
	lex := newLexer(expr)
	var finalCtx map[string]any
	if ctx != nil {
		finalCtx = ctx
	}
	return &parser{lex: lex, peek: lex.next(), ctx: finalCtx}
}

func (p *parser) next() token {
	tok := p.peek
	p.peek = p.lex.next()
	return tok
}

func (p *parser) expect(tt tokenType) token {
	tok := p.next()
	if tok.typ != tt {
		panic(fmt.Sprintf("expected %v, got %v", tt, tok.typ))
	}
	return tok
}

func (p *parser) ParseExpression() bool {
	return p.parseExpression()
}

func (p *parser) parseExpression() bool {
	left := p.parseTerm()
	for p.peek.typ == tOr {
		p.next()
		right := p.parseTerm()
		left = left || right
	}
	return left
}

func (p *parser) parseTerm() bool {
	left := p.parseFactor()
	for p.peek.typ == tAnd {
		p.next()
		right := p.parseFactor()
		left = left && right
	}
	return left
}

func (p *parser) parseFactor() bool {
	if p.peek.typ == tLParen {
		p.next()
		val := p.parseExpression()
		p.expect(tRParen)
		return val
	}
	return p.parseComparison()
}

func (p *parser) parseComparison() bool {
	leftTok := p.next()
	op := p.next()
	rightTok := p.next()

	if op.typ != tOp {
		panic("expected operator")
	}

	leftVal := p.resolveValue(leftTok)
	rightVal := p.resolveValue(rightTok)

	return compare(leftVal, rightVal, op.val)
}

func (p *parser) resolveValue(tok token) any {
	switch tok.typ {
	case tIdent:
		if v, ok := p.ctx[tok.val]; ok {
			return v
		}
		return ""
	case tString:
		return tok.val
	case tNumber:
		f, _ := strconv.ParseFloat(tok.val, 64)
		return f
	}
	return nil
}
