# CrawlCrawl by Kirk Zimmer
Golang Demo BFS Concurrent Webcrawler

## To Build
go build

## To Run

./CrawlCrawl

## To Test (May Require Go 1.17)
go version (to confirm version type)
go test -v


## Program Purpose

1. Fetch the HTML document at the URL
2. Parse out URLs in that HTML document
3. Log/print the URL visited along with all the URLS on the page
4. Loop back to step 1 for each new URL

Also...

1. Maintain a set of all visited URLS to prevent infinite looping / duplicate searching
2. Allow concurrent scraping to multicrawl

## Program Flow

1. Enter in main()
2. Set initial params (startingURL / maxConcurrentWorkers / maxBufferSize)

## Sample Output

![programOutput](https://user-images.githubusercontent.com/8118229/137066075-d2e8def1-0fb0-4ef6-b7b8-21a06142881e.png)
![testOutput](https://user-images.githubusercontent.com/8118229/137066074-9d1e1525-441d-4206-8cae-b81649bbdeff.png)
