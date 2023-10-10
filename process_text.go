package main

import (
	"strings"
	"unicode"
)

func removeInterpunction(input string) string {
	mappingFunc := func(r rune) rune {
		// replace - and _ with spaces; often used as word separators
		if r == '\u002D' || r == '\u005F' {
			return '\u0020'
		}
		//remove other interpunction, but ignore apostrophes because English
		if unicode.IsPunct(r) && r != '\u0027' {
			return -1
		}
		return r
		// return unicode.ToLower(r) // only use lowercase
	}
	return strings.Map(mappingFunc, input)
}

func getWords(data string) []string {
	text := removeInterpunction(data)
	fields := strings.Fields(text) // remove whitespaces

	// data = strings.ReplaceAll(data, "\n", " ")

	return fields
}
