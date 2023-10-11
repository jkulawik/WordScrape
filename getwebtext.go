package main

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// TODO refactor this to return errors
func getWebsiteText(sourceURL string) string {
	response, err := http.Get(sourceURL)
	if err != nil {
		errorLogger.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		warningLogger.Print(response.Status, " -- skipping website: ", sourceURL)
		return ""
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		warningLogger.Print("Unexpected website content type", " -- skipping website: ", sourceURL)
		return ""
	}

	tokenizer := html.NewTokenizer(response.Body)
	previousTokenStartsScript := false
	var fullText string

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err = tokenizer.Err()
			if err == io.EOF {
				break
			} else {
				warningLogger.Print(err)
				break
			}
		} else if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			previousTokenStartsScript = token.Data == "script"
		} else if tokenType == html.TextToken && !previousTokenStartsScript {
			token := tokenizer.Token()
			fullText += token.Data
		}
	}
	return fullText
}
