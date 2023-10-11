package main

import (
	"testing"
)

/* Didn't have time for full testing but needed this one as a sanity check :) */

func TestGetTopFrequentWords(t *testing.T) {
	testdata := []string{"a", "b", "b", "c", "c", "c", "d", "d", "d", "d"}
	expectedResults := []WordCount{
		{"d", 4},
		{"c", 3},
		{"b", 2},
		{"a", 1},
	}
	results := getTopFrequentWords(testdata, 4)

	if len(expectedResults) != len(results) {
		t.Log(results)
		t.Log(len(testdata))
		t.Fatal("Results have incorrect lenght. Expected 4 but got", len(results))
	}

	correctCount := true
	correctOrder := true
	for i := range results {
		entry := results[i]
		expectedEntry := expectedResults[i]

		correctOrder = (entry.Word == expectedEntry.Word) && correctOrder
		correctCount = (entry.Count == expectedEntry.Count) && correctCount
	}

	t.Log("got:      ", results)
	t.Log("expected: ", expectedResults)

	if !correctOrder {
		t.Fatal("words sorted incorrectly")
	}

	if !correctCount {
		t.Fatal("words counted incorrectly")
	}
}
