package main

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func getWebsiteText(sourceURL string) (string, error) {
	response, err := http.Get(sourceURL)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("received status " + response.Status)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", errors.New("unexpected content type")
	}

	tokenizer := html.NewTokenizer(response.Body)
	previousTokenStartsScript := false //text fields can contain JavaScript
	var fullText string

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err = tokenizer.Err()
			if err == io.EOF {
				break
			} else {
				return "", errors.New("HTML parsing error: " + err.Error())
			}
		} else if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			previousTokenStartsScript = token.Data == "script"
		} else if tokenType == html.TextToken && !previousTokenStartsScript {
			token := tokenizer.Token()
			fullText += token.Data
		}
	}
	return fullText, nil
}

func getWordsFromURL(sourceURL string, skipCache bool) []string {
	if isCacheAvailable(sourceURL) && !skipCache {
		infoLogger.Print("Reading cache for ", sourceURL)
		websiteWords, err := readWordCache(sourceURL)
		if err != nil {
			warningLogger.Print(err)
		} else {
			return websiteWords
		}
	}

	infoLogger.Print("Scraping ", sourceURL)
	fullText, err := getWebsiteText(sourceURL)
	if err != nil {
		errorLogger.Print(err, " -- skipping website ", sourceURL)
		return nil
	}

	websiteWords := getWords(fullText)
	infoLogger.Print("Writing cache for ", sourceURL)
	err = writeWordCache(sourceURL, websiteWords)
	if err != nil {
		warningLogger.Print(sourceURL, err)
	}
	return websiteWords
}
