# CrawlCrawl by Kirk Zimmer
Golang Demo Concurrent Webcrawler

## To Build
go build

![buildSequence](https://user-images.githubusercontent.com/8118229/137076378-8fc095b7-c6f3-4bbf-98db-c9883fe6c900.png)


## To Run (With Sample Output)

./CrawlCrawl (runs with http://www.rescale.com as static default)

![programOutput](https://user-images.githubusercontent.com/8118229/137066075-d2e8def1-0fb0-4ef6-b7b8-21a06142881e.png)

./CrawlCrawl --url=https://startingurl.com

![commandlineInput](https://user-images.githubusercontent.com/8118229/137076048-37c6c44e-ab68-486e-ae17-a666b8aca594.png)

## To Test (May Require Go 1.17)

go version (to confirm version type)

go test -v

![testOutput](https://user-images.githubusercontent.com/8118229/137066074-9d1e1525-441d-4206-8cae-b81649bbdeff.png)

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
3. processURL does the meat of the program - getRequest() -> log -> parseBodyForUrls() -> log children -> add new children urls to the queue -> rinse & repeat

