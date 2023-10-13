package main

import (
	"flag"
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

func main() {
	initLoggers()

	exampleURLsCount := flag.Int("e", 0, "How many pages of quotes.toscrape.com examples (max 10) to add to the scrape list")
	workerCount := flag.Int("w", -1, "How many goroutines to use for scraping. Default is one for every URL")
	addBadURLs := flag.Bool("b", false, "Add invalid and 404 URLs to the scrape list")
	skipCache := flag.Bool("s", false, "Do not read cache, always scrape sites")
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("wordscrape [optional flags] [URLs to scrape]...")
		flag.PrintDefaults()
	}
	flag.Parse()

	infoLogger.Print("Starting WordScrape")
	URLs := []string{}

	// Process flags and arguments
	if *exampleURLsCount < 0 {
		*exampleURLsCount = 10
	}

	for i := 0; i < *exampleURLsCount; i++ {
		url := fmt.Sprintf("https://quotes.toscrape.com/page/%d/", i+1)
		URLs = append(URLs, url)
	}

	if *addBadURLs {
		URLs = append(URLs,
			"https://quotes.toscrape.com/doesntexist/",
			"http://definitelydoesnotexist.io/",
			"ftp://definitelydoesnotexist.io/",
			"justsomestring",
		)
	}

	URLs = append(URLs, flag.Args()...)

	if *workerCount < 1 {
		*workerCount = len(URLs)
	}

	// The program proper
	var allWords []string
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

	// Start worker routines
	var workerWaitGroup sync.WaitGroup
	for i := 0; i < *workerCount; i++ {
		workerWaitGroup.Add(1)
		go func() {
			// Consume jobs
			for url := range jobs {
				words := getWordsFromURL(url, *skipCache)
				results <- words
			}
			workerWaitGroup.Done()
		}()
	}

	// Wait for all jobs to be consumed
	workerWaitGroup.Wait()
	close(results)
	resultsWaitGroup.Wait()

	// Display results
	stats := getTopFrequentWords(allWords, 5)
	fmt.Println("\nResults:")
	for _, entry := range stats {
		fmt.Println(entry.Word, "\t| count: ", entry.Count)
	}
}
