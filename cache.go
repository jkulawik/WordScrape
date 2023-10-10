package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const cacheDirectory string = "word-cache/"

type WebsiteWordsCache struct {
	// URLHash string   `json:"url_hash"`

	URL   string   `json:"url"`
	Words []string `json:"words"`
}

func hashString(input string) string {
	h := md5.New()
	io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func writeWordCache(sourceURL string, words []string) error {
	cache := WebsiteWordsCache{
		URL:   sourceURL,
		Words: words,
	}

	cacheJSON, err := json.MarshalIndent(cache, "", "\t")
	if err != nil {
		return errors.New("Error while converting cache to JSON: " + err.Error())
	}

	err = os.MkdirAll(cacheDirectory, 0777)
	if err != nil {
		return errors.New("Error while creating cache directory: " + err.Error())
	}

	filename := hashString(sourceURL)
	err = os.WriteFile(cacheDirectory+filename+".json", cacheJSON, 0666)
	if err != nil {
		return errors.New("Error while writing cache file: " + err.Error())
	}

	return nil
}
