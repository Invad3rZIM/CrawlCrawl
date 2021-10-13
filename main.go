package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("Crawl Crawl by Kirk Zimmer")

	/// 			-------Process-------
	/// 			1. take in input url
	///				2. parse it, add to list of visited

	visitedURLCache := map[string]struct{}{}    //visited url tracking
	concurrencySafetyLock := make(chan bool, 1) //without this, there is synchrony between a map read / write, and crashes program

	startingURL := "http://www.rescale.com"
	maxConcurrentWorkers := 5
	maxURLBufferSize := 1000

	webscrapeQueue := make(chan string, maxURLBufferSize)
	webscrapeQueue <- startingURL

	//enable concurrency
	for i := 0; i < maxConcurrentWorkers; i = i + 1 {
		callback := func() {
			for {
				for len(webscrapeQueue) > 0 {
					processURL(webscrapeQueue, visitedURLCache, concurrencySafetyLock)
				}
			}
		}
		go callback() //go routine callbacks = concurrency requirement :)
	}

	select {} //nonblocking infinite loop.
}

///iterate through
func processURL(webscrapeQueue chan (string), visitedURLCache map[string]struct{}, safetyLock chan (bool)) {
	url := <-webscrapeQueue

	safetyLock <- true
	if _, found := visitedURLCache[url]; found {
		<-safetyLock
		return
	}
	visitedURLCache[url] = struct{}{} //cache result upon visiting
	<-safetyLock

	fmt.Print("\n" + url) //output url

	verbosity := 0
	output, err := getRequest(url)

	if err != nil {
		if verbosity > 0 {
			fmt.Println("\nError: ", url, err.Error())
		}
	} else {
		results := parseBodyForURLs(output)

		for _, childURL := range results {
			safetyLock <- true
			if _, found := visitedURLCache[childURL]; !found {
				<-safetyLock
				webscrapeQueue <- childURL
				fmt.Print("\n\t", childURL) //output children
			} else {
				<-safetyLock
			}
		}
	}
}

//simple get requests, returns the body as a string if no error
func getRequest(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	return string(body), err
}

//parseBodyForURLs is a simple parser that outputs an array of http(s) urls
func parseBodyForURLs(body string) []string {
	currentIteration := 0
	maxIterationAllowed := 99999999 //guard against something relatively close to infinity...

	verbosity := 0 //for debug
	results := []string{}

	subBody := body + "" //no mutation

	for currentIteration < maxIterationAllowed && len(subBody) > 0 {
		startIndex_http := strings.Index(subBody, "http:") // ':' appears to be needed to handle stuff like http-equiv="X-UA-Compatibl....
		startIndex_https := strings.Index(subBody, "https:")
		startIndex_guardedLowest := minPositiveInt(startIndex_http, startIndex_https)
		if startIndex_guardedLowest == -1 {
			break //true exit condition - nothing is found
		}

		subBody = subBody[startIndex_guardedLowest-1:] //go -1 element to get the encapsulation character ie ' / "
		endIndex_toSearch := subBody[0:1]
		subBody = subBody[1:]                                 //move head up by 1 element
		endIndex := strings.Index(subBody, endIndex_toSearch) //search for the matching index

		if endIndex > 0 {
			parsedURL := subBody[0:endIndex]

			matchesSearchRequirement := strings.HasPrefix(parsedURL, "http:") || strings.HasPrefix(parsedURL, "https:") && strings.Contains(parsedURL, ".") //final safety check on core requiremenet... also it seems that all urls have at least 1 '.' symbol
			hasURLDisqualifier := strings.Contains(parsedURL, " ") || strings.Contains(parsedURL, "\n")                                                     //no line spaces or line breaks... simple protection

			if matchesSearchRequirement && !hasURLDisqualifier {
				results = append(results, parsedURL)
			}
		} else {
			if verbosity > 0 {
				fmt.Println("Parse Error - ", endIndex)
			}
		}
		subBody = subBody[1:]
		currentIteration += 1 //emergency exit condition for safety
	}

	return results
}
