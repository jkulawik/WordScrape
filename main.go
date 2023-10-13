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

func main() {
	exampleURLsCount := flag.Int("e", 0, "How many pages of quotes.toscrape.com examples (max 10) to add to the scrape list")
	workerCount := flag.Int("w", 0, "How many goroutines to use for scraping. Default is one for every URL")
	topWordsCount := flag.Int("t", 5, "How many of the most frequently used words to display")
	addBadURLs := flag.Bool("b", false, "Add invalid and 404 URLs to the scrape list")
	skipCache := flag.Bool("s", false, "Do not read cache, always scrape sites")
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("wordscrape [optional flags] [URLs to scrape]...")
		flag.PrintDefaults()
	}
	flag.Parse()

	URLs := []string{}

	/* --- Process flags and arguments --- */
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

	URLs = append(URLs, flag.Args()...) // add remaining non-flag args to URLs

	if *workerCount < 1 {
		*workerCount = len(URLs)
	}

	/* --- The program proper --- */
	initLoggers()
	infoLogger.Print("Starting WordScrape")

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
		resultsWaitGroup.Done()
	}()

	// Start worker routines
	var workerWaitGroup sync.WaitGroup
	for i := 0; i < *workerCount; i++ {
		workerWaitGroup.Add(1)
		go func() {
			// Consume jobs and pipe the results
			for url := range jobs {
				words := getWordsFromURL(url, *skipCache)
				results <- words
			}
			workerWaitGroup.Done()
		}()
	}

	// Wait for all jobs to be finished
	workerWaitGroup.Wait()
	close(results)
	// Wait for the results routine to finish appending all of them
	resultsWaitGroup.Wait()
	/*
		Without waiting for the results, there is a race condition between appending the last result and the last part of the program,
		sometimes ending with the words from the last URL being excluded from the stats.

		We also can't use a single wait group (i.e. *not* wait for the workers to finish before closing the results channel),
		or the results channel will block the program waiting for a result that won't come.
		Alternatively, we could do without workerWaitGroup if we closed the channel by
		only receiving a set number  of results (len(URLs) to be precise), but that feels crude.
	*/

	// Display results
	stats := getTopFrequentWords(allWords, *topWordsCount)
	fmt.Println("\nResults:")
	for _, entry := range stats {
		fmt.Println(entry.Word, "\t| count: ", entry.Count)
	}
}
