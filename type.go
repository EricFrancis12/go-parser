package main

type EnumVariant struct {
	Key   string
	Value string
}

type Enum struct {
	Name     string
	Variants []EnumVariant
}

func (g Enum) gen() string {
	return ""
}

type Struct struct {
	Name   string
	Fields map[string]string
}

func (g Struct) gen() string {
	return ""
}
