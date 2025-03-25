package main

type GenContext struct {
	Fmt string
}

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

type AttributesDirective struct {
	Attributes []Directive
}

func (d AttributesDirective) apply(ctx GenContext, u any) string {
	s := ""
	for _, a := range d.Attributes {
		s += a.apply(ctx, u)
	}
	return s
}

type DeriveAttribute struct {
	TraitNames []TraitName
}

func (d DeriveAttribute) apply(ctx GenContext, unknown any) string {
	s := ""
	for _, tn := range d.TraitNames {
		s += tn.apply(ctx, unknown)
	}
	return s
}

// Trait names found inside of derive()
type TraitName string

const (
	TraitNameClone    TraitName = "Clone"
	TraitNameVariants TraitName = "Variants"
)

func TraitNameFromString(s string) *TraitName {
	tnVariants := []TraitName{
		TraitNameClone,
		TraitNameVariants,
	}

	for _, tn := range tnVariants {
		if string(tn) == s {
			return &tn
		}
	}

	return nil
}

func (d TraitName) apply(ctx GenContext, unknown any) string {
	switch d {
	case TraitNameVariants:
		switch u := unknown.(type) {
		case Enum:
			return fmtEnum(u, ctx)
		}
	}
	return ""
}

type WithDirectives[T any] struct {
	Value      T
	Directives []Directive
}

func (g WithDirectives[T]) gen(ctx GenContext) string {
	s := ""
	for _, d := range g.Directives {
		s += d.apply(ctx, g.Value)
	}
	return s
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
type EnumWithDirectives = WithDirectives[Enum]

// example:
//
// // #[derive(Variants)]
// type Foo struct {}
type StructWithDirectives = WithDirectives[Struct]

func genAll(ctx GenContext, gens ...Generator) string {
	s := ""
	for _, g := range gens {
		s += g.gen(ctx)
	}
	return s
}
