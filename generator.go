package main

import "strings"

type Generator interface {
	gen(GenContext) string
}

type GenContext struct {
	Fmt string
}

// example:
//
// // #[derive(Variants)]
// type Foo struct {}
type Decorated[T any] struct {
	Value      T
	Decorators []Generator
}

func (g Decorated[T]) gen(ctx GenContext) string {
	return genAll(ctx, g.Decorators...)
}

// example:
//
// // #[derive(Variants)]
// type Bar string
//
// const (
//
//	BarOne Bar = "one"
//	BarTwo Bar = "two"
//	BarThree Bar = "three"
//
// )
type DecoratedEnum = Decorated[Enum]

// A type of decorator that contains a list of attributes.
// valid formats:
//
// - // #[...]
// - //#[...]
type AttributesDecorator struct {
	Attributes []Generator
}

func (g AttributesDecorator) gen(ctx GenContext) string {
	return genAll(ctx, g.Attributes...)
}

func genAll(ctx GenContext, gens ...Generator) string {
	var builder strings.Builder
	for _, g := range gens {
		builder.WriteString(g.gen(ctx))
	}
	return builder.String()
}
