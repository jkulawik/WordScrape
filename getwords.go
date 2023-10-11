package main

import (
	"sort"
	"strings"
	"unicode"
)

type WordCount struct {
	Word  string
	Count int
}

func removeInterpunction(input string) string {
	mappingFunc := func(r rune) rune {
		// replace - and _ with spaces; often used as word separators
		if r == '\u002D' || r == '\u005F' {
			return '\u0020'
		}
		//remove interpunction, but ignore apostrophes because English
		if unicode.IsPunct(r) && r != '\u0027' {
			return -1
		}
		return r
	}
	return strings.Map(mappingFunc, input)
}

func getWords(data string) []string {
	text := removeInterpunction(data)
	return strings.Fields(text) // remove whitespaces
}

func getTopFrequentWords(data []string, topCount int) []WordCount {
	wordCountMap := make(map[string]int)
	for _, word := range data {
		lowCaseWord := strings.ToLower(word)
		wordCountMap[lowCaseWord]++
		// this results with a map of unique words and their counts
	}
	/*
		NOTE: maps in Go order their keys independently from the insertion order.
		This means they cannot be used for returning sorted data as-is.
		There is also no map sorting in the standard lib, so a custom solution is needed anyway;
		hence I decided to use a slice of structs instead of map[string]int.

		map[string]int is still used in the word counting above,
		because using a struct slice there would be O(n^2); the current solution is O(n*log n).
		https://stackoverflow.com/questions/29677670/what-is-the-big-o-performance-of-maps-in-golang
	*/
	var wordCounts []WordCount
	for word := range wordCountMap {
		entry := WordCount{word, wordCountMap[word]}
		wordCounts = append(wordCounts, entry)
	}

	sortFunc := func(i, j int) bool {
		return wordCounts[i].Count > wordCounts[j].Count
	}

	sort.Slice(wordCounts, sortFunc)

	realTopCount := 0
	if topCount < 1 || topCount > len(wordCounts) {
		realTopCount = len(wordCounts)
	} else {
		realTopCount = topCount
	}

	return wordCounts[:realTopCount]
}
