package main

type TokenKind int

const (
	EOF TokenKind = iota
	NULL
	TRUE
	FALSE
	NUMBER
	STRING
	IDENTIFIER

	// Grouping & Braces
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	// Equivilance
	ASSIGNMENT
	EQUALS
	NOT_EQUALS
	NOT

	// Conditional
	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	// Logical
	OR
	AND

	// Symbols
	DOT
	DOT_DOT
	SEMI_COLON
	COLON
	QUESTION
	COMMA

	// Shorthand
	PLUS_PLUS
	MINUS_MINUS
	PLUS_EQUALS
	MINUS_EQUALS
	NULLISH_ASSIGNMENT // ??=

	//Maths
	PLUS
	DASH
	SLASH
	STAR
	PERCENT

	// Reserved Keywords
	LET
	CONST
	CLASS
	IMPORT
	FROM
	FN
	IF
	ELSE
	FOREACH
	WHILE
	FOR
	EXPORT
	TYPEOF
	IN

	// Misc
	NUM_TOKENS

	// New
	ENUM
	COMMENT_DIRECTIVE
	TYPE
	IOTA
	DERIVE
	HASHTAG
	UNKNOWN
)

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	"true":    TRUE,
	"false":   FALSE,
	"null":    NULL,
	"let":     LET,
	"const":   CONST,
	"class":   CLASS,
	"import":  IMPORT,
	"from":    FROM,
	"fn":      FN,
	"if":      IF,
	"else":    ELSE,
	"foreach": FOREACH,
	"while":   WHILE,
	"for":     FOR,
	"export":  EXPORT,
	"typeof":  TYPEOF,
	"in":      IN,
	"enum":    ENUM,
	"type":    TYPE,
	"iota":    IOTA,
	"derive":  DERIVE,
}

type Token struct {
	Kind  TokenKind
	Value string
}

func newToken(kind TokenKind, value string) Token {
	return Token{
		kind, value,
	}
}
