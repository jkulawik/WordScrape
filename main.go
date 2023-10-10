package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	warningLogger *log.Logger
	errorLogger   *log.Logger
	infoLogger    *log.Logger
)

func initLoggers() {
	const (
		infoPrefix  = "[INFO]    "
		warnPrefix  = "[WARNING] "
		errorPrefix = "[ERROR]   "
	)

	infoLogger = log.New(os.Stdout, infoPrefix, log.Lshortfile)
	warningLogger = log.New(os.Stderr, warnPrefix, log.Lshortfile)
	errorLogger = log.New(os.Stderr, errorPrefix, log.Lshortfile)
}

func main() {
	initLoggers()
	infoLogger.Print("Starting WordScrape")

	URL := "https://quotes.toscrape.com/page/2/"
	// URL = "https://www.moddb.com/news/an-unfortunate-delay-yet-plenty-of-good-news"

	response, err := http.Get(URL)
	if err != nil {
		errorLogger.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		warningLogger.Print(response.Status, " -- skipping website: ", URL)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		warningLogger.Print("Unexpected website content type", " -- skipping website: ", URL)
	}

	tokenizer := html.NewTokenizer(response.Body)
	var fullText string

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err = tokenizer.Err()
			if err == io.EOF {
				break
			} else {
				errorLogger.Print(err)
				break
			}
		}

		token := tokenizer.Token()
		if tokenType == html.TextToken {
			fullText += token.Data
		}
	}
	words := getWords(fullText)
	fmt.Print(words)
	// fmt.Printf("%q", text) // adds quotes around each element
}
