package main

type EnumVariant struct {
	Key   string
	Value string
}

type Enum struct {
	Name     string
	Variants []EnumVariant
}
