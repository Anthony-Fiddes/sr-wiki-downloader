package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

type WikiPageListingResponse struct {
	Kind string   `json:"kind"`
	Data []string `json:"data"`
}

type WikiPageResponse struct {
	Kind string `json:"kind"`
	Data struct {
		// Some revision data omitted
		ContentMd   string `json:"content_md"`
		Reason      any    `json:"reason"`
		RevisionID  string `json:"revision_id"`
		ContentHTML string `json:"content_html"`
	} `json:"data"`
}

const urlTemplate = "https://www.reddit.com/r/%s/wiki/"

// TODO: take this as a CLI argument
var outputDir = "."

func getPage(subreddit string, page string, outputDir string) error {
	baseURL := fmt.Sprintf(urlTemplate, subreddit)
	pageURL := baseURL + page + ".json"
	response, err := http.Get(pageURL)
	if err != nil {
		return fmt.Errorf("Request to get page data from %s failed: %w", pageURL, err)
	}
	pageResponseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Could not read the response body from %s: %w", pageURL, err)
	}
	pageResponse := WikiPageResponse{}
	err = json.Unmarshal(pageResponseBytes, &pageResponse)
	if err != nil {
		return fmt.Errorf("Could not parse the JSON from %s: %w", pageURL, err)
	}

	markdownBytes := []byte(pageResponse.Data.ContentMd)
	folder := path.Dir(page)
	finalDir := path.Join(outputDir, folder)
	err = os.MkdirAll(finalDir, 0750)
	if err != nil {
		return fmt.Errorf("Could not create the folder '%s': %w", finalDir, err)
	}
	finalFilePath := path.Join(finalDir, path.Base(page)+".md")
	err = os.WriteFile(finalFilePath, markdownBytes, 0644)
	if err != nil {
		return fmt.Errorf("Could not write the markdown file for the wiki page to '%s': %w", finalFilePath, err)
	}

	return nil
}

func main() {
	// TODO: add this from CLI
	subreddit := "germany"
	baseURL := fmt.Sprintf(urlTemplate, subreddit)
	pageListingURL := baseURL + "pages.json"

	pagesResponse, err := http.Get(pageListingURL)
	if err != nil {
		log.Panicf("Request to get list of wiki pages from %s failed: %s", pageListingURL, err)
	}
	pagesBytes, err := io.ReadAll(pagesResponse.Body)
	if err != nil {
		log.Panicf("Could not read the response body from %s: %s", pageListingURL, err)
	}
	wikiPageListing := WikiPageListingResponse{}
	err = json.Unmarshal(pagesBytes, &wikiPageListing)
	if err != nil {
		log.Panicf("Could not parse the list of wiki pages for r/%s: %s", subreddit, err)
	}
	fmt.Println(wikiPageListing)

	for _, page := range wikiPageListing.Data {
		err = getPage(subreddit, page, outputDir)
		if err != nil {
			log.Printf("Could not get r/%s/wiki/%s: %s", subreddit, page, err)
		}
	}
}
