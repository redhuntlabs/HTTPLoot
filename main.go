package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println(fmt.Sprintf(LACKOFART, VERSION))
	flag.IntVar(&httpTimeout, "timeout", 10, "The default timeout for HTTP requests")
	flag.IntVar(&baseFormLen, "form-length", 5, "Length of the string to be randomly generated for filling form fields")
	flag.BoolVar(&submitForm, "submit-forms", false, "Whether to auto-submit forms to trigger debug pages")
	flag.IntVar(&maxWorkers, "concurrency", 100, "Maximum number of sites to process concurrently")
	flag.IntVar(&concurrentURLs, "parallelism", 15, "Number of URLs per site to crawl parallely")
	flag.BoolVar(&wildcardCrawl, "wildcard-crawl", false, "Allow crawling of links outside of the domain being scanned")
	flag.BoolVar(&verifySSL, "verify-ssl", false, "Verify SSL certificates while making HTTP requests")
	flag.IntVar(&crawlDepth, "depth", 3, "Maximum depth limit to traverse while crawling")
	flag.StringVar(&formString, "form-string", "", "Value with which the tool will auto-fill forms, strings will be randomly generated if no value is supplied")
	// flag.IntVar(&MAXCRAWLVAL, "max-crawl", 1000, "Maximum number of links to traverse per site")
	flag.StringVar(&outCsv, "output-file", "httploot-results.csv", "CSV output file path to write the results to")
	flag.StringVar(&userAgent, "user-agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:98.0) Gecko/20100101 Firefox/98.0", "User agent string to use during HTTP requests")
	flag.StringVar(&inpFile, "input-file", "", "Path of the input file containing domains to process")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 && len(inpFile) < 1 {
		log.Fatalln("You need to supply at least a target for the tool to work!")
	}

	_, cancel := context.WithCancel(context.Background())

	writer.Write([]string{"key", "asset", "secret_url", "secret", "stack"})
	if len(dbData) < 1 {
		if err := serializeDB(DATAFILE); err != nil {
			log.Fatalln("Cannot serialize database. exiting...", err.Error())
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go handleInterrupt(c, &cancel)

	tnoe := time.Now()
	log.Println("Starting scan at:", tnoe.Local().String())
	go ProcessHosts(args)

	InitDispatcher(maxWorkers)
	dnoe := time.Now()
	log.Println("Scan finished at:", dnoe.Local().String())
	log.Println("Total time taken:", time.Since(tnoe).String())
	writer.Flush()
}
