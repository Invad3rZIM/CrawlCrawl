package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var inputURL = flag.String("url", "https://legal.twitter.com/imprint", "input url")

func main() {
	flag.Parse() //required to read user input from command line

	fmt.Println("Crawl Crawl by Kirk Zimmer")

	visitedURLCache := map[string]bool{}        //visited url tracking
	concurrencySafetyLock := make(chan bool, 1) //without this, there isn't synchrony between a map read / write, and crashes program

	maxConcurrentWorkers := 5
	maxURLBufferSize := 999999

	webscraperQueue := make(chan string, maxURLBufferSize)
	webscraperQueue <- *inputURL

	//okay...revise this slightly to handle task knowledge better...

	workersStatus := []int{}

	exitSignal := make(chan bool, 1)

	for i := 0; i < maxConcurrentWorkers; i = i + 1 {
		workersStatus = append(workersStatus, 0)
	}
	for i := 0; i < maxConcurrentWorkers; i = i + 1 {
		callback := func(i int) {
			for {
				url := <-webscraperQueue
				workersStatus[i] = 1
				processURL(url, webscraperQueue, visitedURLCache, concurrencySafetyLock)

				allVacant := false

				for j := 0; j < len(workersStatus); j = j + 1 {
					if 1 == workersStatus[j] && j != i {
						allVacant = true
					}
				}
				workersStatus[i] = 0

				time.Sleep(time.Second)

				if allVacant == true && len(webscraperQueue) == 0 {
					exitSignal <- true
				}
			}
		}
		go callback(i)
	}
	for {
		<-exitSignal
		fmt.Println("exiting")
		if len(webscraperQueue) == 0 {
			return
		}
	}
}

///processURL handles a single scrape request, protecting the cache via a lock channel to avoid concurrent read/write crashing
func processURL(url string, webscraperQueue chan (string), visitedURLCache map[string]bool, concurrencySafetyLock chan (bool)) {
	verbosity := 0

	if verbosity > 0 {
		fmt.Println("Q size", len(webscraperQueue))
	}

	concurrencySafetyLock <- true
	if visitedURLCache[url] == true {
		<-concurrencySafetyLock
		return
	}
	<-concurrencySafetyLock

	fmt.Print("\n" + url)       //output url
	visitedURLCache[url] = true //cache result upon visiting

	output, err := getRequest(url)

	if err != nil {
		if verbosity > 0 {
			fmt.Println("\nError: ", url, err.Error())
		}
	} else {
		results := parseBodyForURLs(output)

		for _, childURL := range results {
			concurrencySafetyLock <- true

			// fmt.Println(visitedURLCache, " XXX")
			if visitedURLCache[childURL] == false {
				<-concurrencySafetyLock
				webscraperQueue <- childURL
				fmt.Print("\n\t" + childURL) //output children
			} else {
				<-concurrencySafetyLock
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

	verbosity := 0 //1 for debug

	uniqueResults := map[string]struct{}{}
	results := []string{}

	subBody := body + "" //no mutation

	missingExpectedHREF := false
	missingExpectedHTTP := false
	lastKnownSubBodyBeforeMissingError := ""

	//continuously advance the head on the subBody pointer in O(N^2...good enough TC, but could be optimized if this was truly a bottleneck...)
	//This takes O(N) in a properly structured HTML doc, O(N^2) in a malformed one (ex - an <a tag that didn't contain an href...) let's assume it's O(N^2 since life's imperfect)
	for currentIteration < maxIterationAllowed && len(subBody) > 0 {
		//start with a tags.
		startIndex_aTag := strings.Index(subBody, "<a")
		if startIndex_aTag == -1 {
			break //loop exit condition
		}
		subBody = subBody[startIndex_aTag:]

		//find the closest href
		startIndex_href := strings.Index(subBody, "href")
		if startIndex_href == -1 {
			missingExpectedHREF = true
			lastKnownSubBodyBeforeMissingError = subBody
			break //loop exit condition
		}
		subBody = subBody[startIndex_href:]

		//find the closest starting token after the href
		startIndex_http := strings.Index(subBody, "http:") // ':' appears to be needed to handle stuff like http-equiv="X-UA-Compatibl....
		startIndex_https := strings.Index(subBody, "https:")
		startIndex_guardedLowest := minPositiveInt(startIndex_http, startIndex_https)

		if startIndex_guardedLowest == -1 {
			missingExpectedHTTP = true
			lastKnownSubBodyBeforeMissingError = subBody
			break //loop exit condition
		}

		//find the appropriate ending token
		subBody = subBody[startIndex_guardedLowest-1:] //		go -1 element to get the encapsulation character ie ' / "
		endToken_toSearch := subBody[0:1]
		subBody = subBody[1:]                                 //move head up by 1 element
		endIndex := strings.Index(subBody, endToken_toSearch) //search for the matching index

		if endIndex > 0 {
			parsedURL := subBody[0:endIndex]

			matchesSearchRequirement := strings.HasPrefix(parsedURL, "http:") || strings.HasPrefix(parsedURL, "https:") && strings.Contains(parsedURL, ".") //final safety check on core requiremenet... also it seems that all urls have at least 1 '.' symbol
			hasURLDisqualifier := strings.Contains(parsedURL, " ") || strings.Contains(parsedURL, "\n")                                                     //no line spaces or line breaks... simple protection

			if matchesSearchRequirement && !hasURLDisqualifier {
				trimmedURL := strings.Trim(parsedURL, " ")
				uniqueResults[trimmedURL] = struct{}{}
			} else {
				if verbosity > 0 {
					fmt.Println("Discarding result : ", parsedURL, " after checks : on matchesSearchRequirement, hasURLDisqualifier ( ", matchesSearchRequirement, hasURLDisqualifier, " ) should ideally be ( true, false )")
				}
			}
		} else {
			if verbosity > 0 {
				fmt.Println("Parse Error - ", endIndex)
			}
		}
		subBody = subBody[1:] //advance pointer up 1, since in this phase it points to ->[h]ttp...
		currentIteration += 1 //emergency exit condition for safety
	}

	//for debugging information on the exit condition that terminated the loop
	if verbosity > 0 {
		if missingExpectedHREF {
			fmt.Println("Missing Expected HREF tag: (parse break location : ", lastKnownSubBodyBeforeMissingError, ")")
		}
		if missingExpectedHTTP {
			fmt.Println("Missing Expected HTTP text : (parse break location : ", lastKnownSubBodyBeforeMissingError, ")")
		}
		if len(subBody) == 0 {
			fmt.Println("Parse loop exited after processing entire body")
		}
		fmt.Println("Parse loop exited after iteration count : ", currentIteration)
	}

	for key, _ := range uniqueResults {
		results = append(results, key)
	}

	return results
}
