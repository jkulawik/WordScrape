package main

import (
	"fmt"
	"log"
	"os"
	"sync"
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
		errorLogger.Print(sourceURL, " ", err, " -- website will be skipped")
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

func getWordsFromURLWorker(jobs chan string, results chan []string, waitGroup *sync.WaitGroup) {
	// Consume jobs
	for url := range jobs {
		words := getWordsFromURL(url)
		results <- words
	}
	waitGroup.Done()
}

func main() {
	initLoggers()
	infoLogger.Print("Starting WordScrape")

	URLs := []string{
		// "https://quotes.toscrape.com/doesntexist/",
		// "https://www.moddb.com/news/doesntexist/",
		"https://quotes.toscrape.com/page/2/",
		"https://quotes.toscrape.com/page/3/",
		"https://quotes.toscrape.com/page/4/",
		"https://quotes.toscrape.com/page/5/",
		// "https://www.moddb.com/news/an-unfortunate-delay-yet-plenty-of-good-news",
	}

	var allWords []string
	workerCount := 2
	jobs := make(chan string)
	results := make(chan []string)

	// Start filling the job pipeline
	go func() {
		for _, url := range URLs {
			jobs <- url
		}
		close(jobs)
	}()

	// Start consuming results
	var resultsWaitGroup sync.WaitGroup
	resultsWaitGroup.Add(1)
	go func() {
		for result := range results {
			allWords = append(allWords, result...)
		}
		// must wait for the last append to finish;
		// without it there was a race condition where the last result was sometimes not being included in the stats
		resultsWaitGroup.Done()
	}()

	// Start workers
	var workerWaitGroup sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		workerWaitGroup.Add(1)
		go getWordsFromURLWorker(jobs, results, &workerWaitGroup)
	}

	// Wait for all jobs to be consumed
	workerWaitGroup.Wait()
	close(results)
	resultsWaitGroup.Wait()

	stats := getTopFrequentWords(allWords, 5)
	fmt.Println("\nResults:")
	for _, entry := range stats {
		fmt.Println(entry.Word, "\t| count: ", entry.Count)
	}
}
