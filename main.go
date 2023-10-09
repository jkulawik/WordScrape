package main

import (
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
	infoLogger.Print("WordScrape")

	URL := "https://quotes.toscrape.com/page/2/aaaa"

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
}
