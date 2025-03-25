package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var lookupFuncs = []LookupFunc[Generator]{
	// The parser starts at index 0
	// and moves up until it finds a match,
	// so naturally items with a lower index,
	// have higher priority of being matched.
	parseDecorator,
}

func parseDecorator(p *Parser[Generator]) (Generator, bool) {
	startingPos := p.GetPos()

	fns := []LookupFunc[Generator]{
		parseAttributesDecorator,
	}

	for _, fn := range fns {
		if gen, ok := fn(p); ok {
			return gen, true
		}

		p.SetPos(startingPos)
	}

	return nil, false
}

// Attempts to parse:
// // #[...] or //#[...]
//
// Returns a *AttributesDecorator{}
func parseAttributesDecorator(p *Parser[Generator]) (Generator, bool) {
	return &AttributesDecorator{}, true
}

func TestParseCommentDirectives(t *testing.T) {
	type Test struct {
		source string
	}

	tests := []Test{
		{
			source: "//#[derive(Variants)]",
		},
		{
			source: "// #[derive(Variants)]",
		},
	}

	for _, test := range tests {
		tokens := Tokenize(test.source)
		p := NewParser(tokens, lookupFuncs)

		assert.Equal(t, COMMENT_DIRECTIVE, p.CurrentTokenKind())
		assert.Len(t, p.tokens, 2)
		assert.Equal(t, COMMENT_DIRECTIVE, p.tokens[0].Kind)
		assert.Equal(t, EOF, p.tokens[1].Kind)

		gen, ok := p.Match()
		assert.True(t, ok)
		assert.Equal(t, "", gen.gen(GenContext{}))

		ad, ok := interface{}(gen).(*AttributesDecorator)
		assert.True(t, ok)
		assert.Len(t, ad.Attributes, 0)
	}
}
