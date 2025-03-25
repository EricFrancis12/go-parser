package main

type EnumVariant struct {
	Key   string
	Value string
}

type Enum struct {
	Name     string
	Variants []EnumVariant
}

func (g Enum) gen(GenContext) string {
	return ""
}

type Struct struct {
	Name   string
	Fields map[string]string
}

func (g Struct) gen(GenContext) string {
	return ""
}
