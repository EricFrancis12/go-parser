package main

import "strings"

func GoCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		newWord := ""
		for j, char := range word {
			if j == 0 {
				newWord += strings.ToUpper(string(char))
			} else {
				newWord += strings.ToLower(string(char))
			}
		}
		words[i] = newWord
	}
	return strings.Join(words, "")
}
