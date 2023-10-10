package main

import (
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

func main() {
	initLoggers()
	infoLogger.Print("Starting WordScrape")

	URL := "https://quotes.toscrape.com/page/2/"
	// URL = "https://www.moddb.com/news/an-unfortunate-delay-yet-plenty-of-good-news"

	fullText := getWebsiteText(URL)
	words := getWords(fullText)
	// fmt.Print(words)
	// fmt.Printf("%q", text) // adds quotes around each element

	err := writeWordCache(URL, words)
	if err != nil {
		warningLogger.Print(err)
	}
}
