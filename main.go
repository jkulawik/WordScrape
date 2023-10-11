package main

import (
	"fmt"
	"log"
	"os"
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

func getWordsFromURL(sourceURL string) []string {
	if isCacheAvailable(sourceURL) {
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
		warningLogger.Print(sourceURL, err, " -- website words will be skipped")
		return nil
	}

	websiteWords := getWords(fullText)
	err = writeWordCache(sourceURL, websiteWords)
	if err != nil {
		warningLogger.Print(sourceURL, err)
	}
	return websiteWords
}

func main() {
	initLoggers()
	infoLogger.Print("Starting WordScrape")

	URL := "https://quotes.toscrape.com/page/2/"
	// URL = "https://www.moddb.com/news/an-unfortunate-delay-yet-plenty-of-good-news"

	words := getWordsFromURL(URL)
	fmt.Print(words)
	// fmt.Printf("%q", text) // adds quotes around each element

}
