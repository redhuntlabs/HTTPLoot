package main

import "regexp"

type (
	FoundIssue struct {
		Issue  string `json:"issue"`
		Path   string `json:"path"`
		Type   string `json:"type"`
		Secret string `json:"secret"`
	}
	DBData struct {
		Issue      string   `json:"issue"`
		Severity   string   `json:"severity"`
		Detectors  []string `json:"detectors"`
		Validators struct {
			Status []int    `json:"status"`
			Regex  []string `json:"regex"`
		} `json:"validators"`
		Extractors []struct {
			Regex   string `json:"regex"`
			Cgroups string `json:"cgroups"`
		} `json:"extractors"`
	}
)

var (
	httpTimeout, baseFormLen               int
	maxWorkers, concurrentURLs, crawlDepth int
	wildcardCrawl, submitForm, verifySSL   bool
	inpFile, userAgent, formString, outCsv string

	dbData    map[string]DBData
	regexData map[string]string

	crawlProgress = 0
	reJSScript    = regexp.MustCompile(`(?i)<script[^>]+src=['"]?([^'"\s>]+)`)

	FORM_STRING = "httpl00t"
	VERSION     = "0.1"
	DATAFILE    = "lootdb.json"
	OUTCSV      = "httploot-results.csv"
	REGEXFILE   = "regexes.json"
	USERAGENT   = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:98.0) Gecko/20100101 Firefox/98.0"
	BYTES       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LACKOFART   = `
      _____
       )=(
      /   \     H T T P L O O T
     (  $  )                  v%s
      \___/

[+] Log4jHunt by RedHunt Labs - A Modern Attack Surface (ASM) Management Company
[+] Author: Pinaki Mondal (RHL Research Team)
[+] Continuously Track Your Attack Surface using https://redhuntlabs.com/nvadr.
`
	// MAXCRAWLVAL int
)
