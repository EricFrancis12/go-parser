package main

import (
	"fmt"
	"strings"
)

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

func fmtEnum(enum Enum, ctx GenContext) string {
	var (
		goCase         = false
		enumNamePrefix = false
		namePrefix     = ""
		variantPrefix  = ""
	)
	if ctx.Fmt == "PRISMA" {
		goCase = true
		enumNamePrefix = true
		namePrefix = "db."
		variantPrefix = "db."
	}

	variants := []string{}
	for _, v := range enum.Variants {
		key := v.Key
		if goCase {
			key = GoCase(key)
		}

		npre := ""
		if enumNamePrefix {
			npre = enum.Name
		}

		variants = append(variants, fmt.Sprintf("%s%s%s,", namePrefix, npre, key))
	}

	return fmt.Sprintf(`
			var %sVariants = [%d]%s%s{
				%s
			}
		`,
		enum.Name,
		len(enum.Variants),
		variantPrefix,
		enum.Name,
		strings.Join(variants, "\n"),
	)
}
