package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
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

const (
	urlTemplate      = "https://www.reddit.com/r/%s/wiki/"
	wikiPageTemplate = "r/%s/wiki/%s"
	usageTemplate    = `Usage: %s <subreddit_name> [output_dir]

output_dir defaults to the current diretory if not provided.`
)

var (
	outputDir    = "."
	RateLimitErr = errors.New("reddit rate limited the request")
)

func getPage(subreddit string, page string, outputDir string) error {
	baseURL := fmt.Sprintf(urlTemplate, subreddit)
	pageURL := baseURL + page + ".json"
	response, err := http.Get(pageURL)
	if err != nil {
		return fmt.Errorf("request to %s failed: %w", pageURL, err)
	}
	if response.StatusCode == 429 {
		return fmt.Errorf("request to %s failed: %w", pageURL, RateLimitErr)
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("request to %s failed with error code %s", pageURL, response.Status)
	}

	pageResponseBytes, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return fmt.Errorf("could not read the response body from %s: %w", pageURL, err)
	}
	pageResponse := WikiPageResponse{}
	err = json.Unmarshal(pageResponseBytes, &pageResponse)
	if err != nil {
		return fmt.Errorf("could not parse the JSON from %s: %w", pageURL, err)
	}

	markdownBytes := []byte(pageResponse.Data.ContentMd)
	if len(markdownBytes) == 0 {
		return fmt.Errorf("no markdown content found for r/%s/wiki/%s", subreddit, page)
	}

	folder := path.Dir(page)
	finalDir := path.Join(outputDir, folder)
	err = os.MkdirAll(finalDir, 0750)
	if err != nil {
		return fmt.Errorf("could not create the folder '%s': %w", finalDir, err)
	}
	finalFilePath := path.Join(finalDir, path.Base(page)+".md")
	err = os.WriteFile(finalFilePath, markdownBytes, 0644)
	if err != nil {
		return fmt.Errorf("could not write the markdown file for the wiki page to '%s': %w", finalFilePath, err)
	}

	return nil
}

func printUsageAndFail() {
	log.SetFlags(0)
	log.Println()
	log.Fatalf(usageTemplate, path.Base(os.Args[0]))
}

func main() {
	if len(os.Args) < 2 {
		log.Println("One argument required: subreddit_name")
		printUsageAndFail()
	}
	var subreddit string
	if len(os.Args) >= 2 {
		subreddit = os.Args[1]
	}
	if len(os.Args) == 3 {
		outputDir = os.Args[2]
	}
	if len(os.Args) > 3 {
		log.Printf("%s doesn't take more than 2 arguments", os.Args[0])
		printUsageAndFail()
	}
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

	listingWithIndex := make([]string, 0, len(wikiPageListing.Data)+1)
	listingWithIndex = append(listingWithIndex, "index")
	listingWithIndex = append(listingWithIndex, wikiPageListing.Data...)
	for _, page := range listingWithIndex {
		wikiPage := fmt.Sprintf(wikiPageTemplate, subreddit, page)
		for i := 0; i < 3; i++ {
			retry := false
			err = getPage(subreddit, page, outputDir)
			if errors.Is(err, RateLimitErr) {
				log.Printf("Attempt %d to request %s was rate limited", i, wikiPage)
				retry = true
			} else if err != nil {
				log.Printf("Could not get %s: %s", wikiPage, err)
			} else {
				log.Printf("Successfully downloaded %s", wikiPage)
			}

			// Without authenticating to the Reddit API, you're only allowed 10
			// requests per minute.
			time.Sleep(time.Minute/10 + 1)
			if !retry {
				break
			}
		}
	}
}
