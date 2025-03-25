package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommentDirectives(t *testing.T) {
	type Test struct {
		source string
	}

	tests := []Test{
		{
			source: `
				// #[derive(Variants)]
				type Bar string

				const (
					BarOne Bar = "one"
					BarTwo Bar = "two"
					BarThree Bar = "three"
				)
			`,
		},
		{
			source: `
				// #[derive(Clone, Variants)]
				type Bar = string

				const (
					BarOne Bar = "one"
					BarTwo Bar = "two"
					BarThree Bar = "three"
				)
			`,
		},
		{
			source: `
				// #[derive(Variants)]
				enum DmgType {
					BLUNT
					SLASH
					PIERCE
					FIRE
					WATER
					EARTH
					WIND
					ICE
					GRASS
					ELECTRIC
					ARCANE
					BLOOD
					POISON
					CORROSIVE
					LIGHT
					DARK
				}
			`,
		},
	}

	for _, test := range tests {
		tokens := Tokenize(test.source)
		p := NewParser(tokens, lookupFuncs)

		gen, ok := p.Match()
		assert.True(t, ok)
		assert.Equal(t, "---enum---", gen.gen(GenContext{}))

		ewd, ok := gen.(EnumWithDirectives)
		assert.True(t, ok)
		assert.Len(t, ewd.Directives, 1)
	}
}
