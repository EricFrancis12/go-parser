package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommentDirectives(t *testing.T) {
	type Test struct {
		source   string
		ctx      GenContext
		expected Enum
	}

	tests := []Test{
		{
			source: `
				// #[derive(Variants)]
				type Foo string

				const FooOne Foo = "ONE"
			`,
			expected: Enum{
				Name:     "Foo",
				Variants: []EnumVariant{{"FooOne", "ONE"}},
			},
		},
		{
			source: `
				// #[derive(Variants)]
				type Foo string

				const (
					FooOne Foo = "ONE"
					FooTwo Foo = "TWO"
					FooThree Foo = "THREE"
				)
			`,
			expected: Enum{
				Name:     "Foo",
				Variants: []EnumVariant{{"FooOne", "ONE"}, {"FooTwo", "TWO"}, {"FooThree", "THREE"}},
			},
		},
		{
			source: `
				// #[derive(Clone, Variants)]
				type Foo = string

				const (
					FooOne Foo = "ONE"
					FooTwo Foo = "TWO"
					FooThree Foo = "THREE"
				)
			`,
			expected: Enum{
				Name:     "Foo",
				Variants: []EnumVariant{{"FooOne", "ONE"}, {"FooTwo", "TWO"}, {"FooThree", "THREE"}},
			},
		},
		{
			source: `
				// #[derive(Variants)]
				enum Foo {
					ONE
					TWO
					THREE
				}
			`,
			ctx: GenContext{Fmt: "PRISMA"},
			expected: Enum{
				Name:     "Foo",
				Variants: []EnumVariant{{"ONE", "ONE"}, {"TWO", "TWO"}, {"THREE", "THREE"}},
			},
		},
	}

	for _, test := range tests {
		tokens := Tokenize(test.source)
		p := NewParser(tokens, lookupFuncs)

		gen, ok := p.Match()
		assert.True(t, ok)

		ewd, ok := gen.(EnumWithDirectives)
		assert.True(t, ok)

		assert.Len(t, ewd.Directives, 1)
		assert.Len(t, ewd.Value.Variants, len(test.expected.Variants))

		assert.Equal(t, test.expected.Name, ewd.Value.Name)
		for i, v := range test.expected.Variants {
			assert.Equal(t, v.Key, ewd.Value.Variants[i].Key)
		}
	}
}
