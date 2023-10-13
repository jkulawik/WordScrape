# WordScrape

A programming challenge. Features:
- Scrape websites for text and return the most frequently used words
- Cache results per website to not scrape unnecessarily
- Scrape several websites in parallel with goroutines

## Usage
```
wordscrape [optional flags] [URLs to scrape]...
  -b    Add invalid and 404 URLs to the scrape list
  -e int
        How many pages of quotes.toscrape.com examples (max 10) to add to the scrape list
  -s    Do not read cache, always scrape sites
  -t int
        How many of the most frequently used words to display (default 5)
  -w int
        How many goroutines to use for scraping. Default is one for every URL (default -1)
```

Cache is written to a `word-cache` folder in the directory the program was ran in.

## Examples
Replace `go run .` with the exec name (`wordscrape` by default) if you prefer to compile the program manually (`go build`).

Scrape 10 example pages without reading cache:
```
go run . -e 10 -s
```

Scrape 10 example pages with cache (if cache is available):
```
go run . -e 10
```

Scrape 2 example pages and include invalid URLs in scrape list
```
go run . -e 2 -b
```

Scrape 2 example pages and a list of custom pages:
```
go run . -e 2 https://go.dev/blog/loopvar-preview https://go.dev/blog/wasi
```
