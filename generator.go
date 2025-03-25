package main

import "strings"

type GenContext struct{}

// Generators can be any "trigger" in the source code
// that results in code being generated.
type Generator interface {
	gen(GenContext) string
}

// A directive is an instruction found in the source code
// outlining how Generators are created.
type Directive interface {
	apply(GenContext, any) string
}

type AttributeDirective struct {
	Attributes []Attribute
}

func (d AttributeDirective) apply(ctx GenContext, u any) string {
	switch u.(type) {
	case Enum:
		// TODO: ...
		return "---enum---"
	case Struct:
		// TODO: ...
	}
	return ""
}

type Attribute interface {
	// TODO: ...
}

type DeriveAttribute struct {
	TraitNames []string // TODO: change to TraitName
}

func (d DeriveAttribute) apply(ctx GenContext, u any) string {
	switch u.(type) {
	case Enum:
		// TODO: ...
	case Struct:
		// TODO: ...
	}
	return ""
}

// Trait names found inside of derive()
type TraitName string

const (
	TraitNameClone    TraitName = "Clone"
	TraitNameVariants TraitName = "Variants"
)

type WithDirectives[T any] struct {
	Value      T
	Directives []Directive
}

func (g WithDirectives[T]) gen(ctx GenContext) string {
	builder := strings.Builder{}
	for _, d := range g.Directives {
		builder.WriteString(d.apply(ctx, g.Value))
	}
	return builder.String()
}

// example:
//
// // #[derive(Variants)]
// type Bar string
//
// const (
//	BarOne Bar = "one"
//	BarTwo Bar = "two"
//	BarThree Bar = "three"
// )
type EnumWithDirectives = WithDirectives[Enum]

// example:
//
// // #[derive(Variants)]
// type Foo struct {}
type StructWithDirectives = WithDirectives[Struct]
