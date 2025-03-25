package main

import "fmt"

type LookupFunc[T any] func(p *Parser[T]) (T, bool)

type Parser[T any] struct {
	tokens      []Token
	pos         int
	lookupFuncs []LookupFunc[T]
}

func NewParser[T any](tokens []Token, fns []LookupFunc[T]) *Parser[T] {
	p := &Parser[T]{
		tokens:      tokens,
		pos:         0,
		lookupFuncs: fns,
	}
	return p
}

func (p *Parser[T]) Match() (T, bool) {
	startingPos := p.GetPos()
	for _, fn := range p.lookupFuncs {
		if t, ok := fn(p); ok {
			return t, true
		}

		p.SetPos(startingPos)
	}

	var t T
	return t, false
}

func (p *Parser[T]) MustParseAll() []T {
	result := []T{}

	for p.Advance().Kind != EOF {
		if gs, ok := p.Match(); ok {
			result = append(result, gs)
		}
	}

	return result
}

func (p *Parser[T]) CurrentToken() Token {
	return p.tokens[p.pos]
}

func (p *Parser[T]) Advance() Token {
	tk := p.CurrentToken()
	p.pos++
	return tk
}

func (p *Parser[T]) hasTokens() bool {
	return p.pos < len(p.tokens) && p.CurrentTokenKind() != EOF
}

func (p *Parser[T]) PreviousToken() Token {
	return p.tokens[p.pos-1]
}

func (p *Parser[T]) CurrentTokenKind() TokenKind {
	return p.tokens[p.pos].Kind
}

func (p *Parser[T]) GetPos() int {
	return p.pos
}

func (p *Parser[T]) SetPos(pos int) {
	p.pos = pos
}

func (p *Parser[T]) Reset(pos int) (T, bool) {
	p.SetPos(pos)
	var t T
	return t, false
}

func (p *Parser[T]) expectError(expectedKind TokenKind, err any) Token {
	token := p.CurrentToken()
	kind := token.Kind

	if kind != expectedKind {
		if err == nil {
			err = fmt.Sprintf("Expected %d but recieved %d instead\n", expectedKind, kind)
		}

		panic(err)
	}

	return p.Advance()
}

func (p *Parser[T]) expect(expectedKind TokenKind) Token {
	return p.expectError(expectedKind, nil)
}
