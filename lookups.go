package main

import (
	"fmt"
	"strings"
)

var lookupFuncs = []LookupFunc[Generator]{
	// The parser starts at index 0
	// and moves up until it finds a match,
	// so naturally items with a lower index,
	// have higher priority of being matched.
	_parseCommentDirective,
	_parseEnum,
}

// Returns an EnumWithDirectives{} or StructWithDirectives{}
func _parseCommentDirective(p *Parser[Generator]) (Generator, bool) {
	tk := p.Advance()
	if tk.Kind != COMMENT_DIRECTIVE {
		return nil, false
	}

	if tk.Value[0:2] != "//" {
		panic("expected 2 forward slashes at the beginning of comments")
	}

	// Create new parser to handle tokens inside of comment directive
	tokens := Tokenize(tk.Value[2:])
	_p := NewParser(tokens, []LookupFunc[Directive]{})

	d, ok := parseDirective(_p)
	if !ok {
		return nil, false
	}

	switch d.(type) {
	case AttributesDirective:
		gen, ok := p.Match()
		if ok {
			switch g := gen.(type) {
			case Enum:
				return EnumWithDirectives{
					Value:      g,
					Directives: []Directive{d},
				}, true
			default:
				panic(fmt.Sprintf("unknown generator: %+v\n", gen))
			}
		}
	default:
		panic(fmt.Sprintf("unknown directive: %+v\n", d))
	}

	return nil, false
}

func _parseEnum(p *Parser[Generator]) (Generator, bool) {
	fns := []LookupFunc[Generator]{
		_parseGoEnum,
		_parsePrismaEnum,
	}

	startingPos := p.GetPos()
	for _, fn := range fns {
		if gen, ok := fn(p); ok {
			return gen, true
		}

		p.SetPos(startingPos)
	}

	return nil, false
}

func _parseGoEnum(p *Parser[Generator]) (Generator, bool) {
	startingPos := p.GetPos()

	if p.Advance().Kind != TYPE {
		return p.Reset(startingPos)
	}

	enum := Enum{
		Name:     p.CurrentToken().Value,
		Variants: []EnumVariant{},
	}

	if p.Advance().Kind != IDENTIFIER {
		return p.Reset(startingPos)
	}
	if p.Advance().Kind == ASSIGNMENT {
		p.Advance()
	}

	if p.Advance().Kind != CONST {
		return p.Reset(startingPos)
	}
	if p.CurrentTokenKind() == OPEN_PAREN {
		p.Advance()
	}

	for p.CurrentTokenKind() == IDENTIFIER {
		enumVariant := EnumVariant{
			Key: p.Advance().Value,
		}

		if p.CurrentTokenKind() != IDENTIFIER || p.Advance().Value != enum.Name {
			return p.Reset(startingPos)
		}

		if p.Advance().Kind != ASSIGNMENT {
			return p.Reset(startingPos)
		}

		if p.CurrentTokenKind() == IOTA {
			enumVariant.Value = "0"
			enum.Variants = append(enum.Variants, enumVariant)

			p.Advance()

			for i := 1; p.CurrentTokenKind() == IDENTIFIER; i++ {
				enum.Variants = append(enum.Variants, EnumVariant{
					Key:   p.Advance().Value,
					Value: fmt.Sprintf("%d", i),
				})
			}

			break
		}

		enumVariant.Value = strings.Trim(p.Advance().Value, `"`)
		enum.Variants = append(enum.Variants, enumVariant)
	}

	if p.CurrentTokenKind() != CLOSE_PAREN {
		p.Advance()
	}

	return enum, true
}

func _parsePrismaEnum(p *Parser[Generator]) (Generator, bool) {
	startingPos := p.GetPos()

	if p.Advance().Kind != ENUM {
		return p.Reset(startingPos)
	}

	enum := Enum{
		Name:     p.CurrentToken().Value,
		Variants: []EnumVariant{},
	}

	if _, ok := p.AdvanceTo(IDENTIFIER, OPEN_CURLY); !ok {
		return p.Reset(startingPos)
	}

	for p.CurrentTokenKind() == IDENTIFIER {
		variantName := p.Advance().Value
		enum.Variants = append(enum.Variants, EnumVariant{
			Key:   variantName,
			Value: variantName,
		})
	}

	if p.CurrentTokenKind() != CLOSE_CURLY {
		return p.Reset(startingPos)
	}

	return enum, true
}

func parseDirective(p *Parser[Directive]) (Directive, bool) {
	startingPos := p.GetPos()

	fns := []LookupFunc[Directive]{
		parseAttributesDirective,
	}

	for _, fn := range fns {
		if d, ok := fn(p); ok {
			return d, true
		}

		p.SetPos(startingPos)
	}

	return nil, false
}

func parseAttributesDirective(p *Parser[Directive]) (Directive, bool) {
	startingPos := p.GetPos()

	if p.Advance().Kind != HASHTAG {
		return p.Reset(startingPos)
	}
	if p.Advance().Kind != OPEN_BRACKET {
		return p.Reset(startingPos)
	}

	attributes := []Directive{}

	for p.CurrentTokenKind() != CLOSE_PAREN {
		att, ok := parseAttribute(p)
		if !ok {
			return p.Reset(startingPos)
		}

		if p.CurrentTokenKind() == COMMA {
			p.Advance()
		}

		attributes = append(attributes, att)
	}

	return AttributesDirective{
		Attributes: attributes,
	}, true
}

func parseAttribute(p *Parser[Directive]) (Directive, bool) {
	startingPos := p.GetPos()

	fns := []LookupFunc[Directive]{
		parseDeriveAttribute,
	}

	for _, fn := range fns {
		if d, ok := fn(p); ok {
			return d, true
		}

		p.SetPos(startingPos)
	}

	return nil, false
}

func parseDeriveAttribute(p *Parser[Directive]) (Directive, bool) {
	startingPos := p.GetPos()

	if p.Advance().Kind != DERIVE {
		return p.Reset(startingPos)
	}
	if p.Advance().Kind != OPEN_PAREN {
		return p.Reset(startingPos)
	}

	traitNames := []TraitName{}

	for p.CurrentTokenKind() != CLOSE_PAREN {
		tk := p.Advance()
		if tk.Kind != IDENTIFIER {
			return p.Reset(startingPos)
		}

		if p.CurrentTokenKind() == COMMA {
			p.Advance()
		}

		tn := TraitNameFromString(tk.Value)
		if tn == nil {
			panic(fmt.Sprintf("expected valid Trait, but got: %s\n", tk.Value))
		}

		traitNames = append(traitNames, *tn)
	}

	return DeriveAttribute{
		TraitNames: traitNames,
	}, true
}
