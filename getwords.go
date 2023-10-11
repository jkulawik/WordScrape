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
	return strings.Fields(text) // remove whitespaces
}

func getTopFrequentWords(data []string, topCount int) []WordCount {
	wordCountMap := make(map[string]int)
	for _, word := range data {
		lowCaseWord := strings.ToLower(word)
		wordCountMap[lowCaseWord]++
		// this results with a map of unique words and their counts
	}
	// There is no native map sorting; instead sort the keys in a slice using the map values

	var uniqueWords []string
	for word := range wordCountMap {
		uniqueWords = append(uniqueWords, word)
	}

	sortFunc := func(i, j int) bool {
		return wordCountMap[uniqueWords[i]] > wordCountMap[uniqueWords[j]]
	}

	sort.Slice(uniqueWords, sortFunc)

	realTopCount := 0
	if topCount < 1 || topCount > len(uniqueWords) {
		realTopCount = len(uniqueWords)
	} else {
		realTopCount = topCount
	}

	/*
		NOTE: maps in Go order their keys deterministically, independently from the insertion order.
		This means that regular maps cannot be used for returning sorted data as-is;
		hence I decided to use and return a slice of structs instead of map[string]int.
	*/

	topWords := make([]WordCount, realTopCount)
	for i := range topWords {
		currentKey := uniqueWords[i]
		topWords[i].Word = currentKey
		topWords[i].Count = wordCountMap[currentKey]
	}

	return topWords
}
