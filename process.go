package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

func (p *ProcJob) FingerPrint() {
	log.Println("Trying to identify tech stack:", p.Host)
	heads, body, err := makeRequest(p.Host)
	if err != nil {
		log.Println("Error during tech stack detection:", err.Error())
		return
	}
	fresp := heads + "\n\n" + string(*body)
outBreak:
	for tech, ddata := range dbData {
		for _, drex := range ddata.Detectors {
			rex := regexp.MustCompile(drex)
			if rex.MatchString(fresp) {
				log.Println("Identified tech stack:", tech)
				p.Host = p.Host + ":::" + tech
				break outBreak
			}
		}
	}
}

func (p *ProcJob) ExecuteCrawler() {
	hastechs := false
	mainurl := p.Host

	if strings.Contains(p.Host, ":::") {
		hastechs = true
		mainurl = strings.Split(p.Host, ":::")[0]
	}
	log.Println("Processing:", p.Host)

	if !strings.Contains(p.Host, "://") {
		mainurl = "http://" + p.Host
	} else {
		p.Host = strings.Split(p.Host, "://")[1]
	}

	var c *colly.Collector
	if !WILDCARD_CRAWL {
		c = colly.NewCollector(
			colly.AllowedDomains(strings.Split(mainurl, "://")[1]),
			colly.UserAgent(USERAGENT),
			colly.ParseHTTPErrorResponse(),
			colly.MaxDepth(CRAWL_DEPTH),
			colly.CacheDir(".colly_cache/"),
		)
	} else {
		c = colly.NewCollector(
			colly.UserAgent(USERAGENT),
			colly.ParseHTTPErrorResponse(),
			colly.MaxDepth(CRAWL_DEPTH),
			colly.CacheDir(".colly_cache/"),
		)
	}

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !VERIFY_SSL,
		},
		DisableKeepAlives: true, // we disable keep alive targets
	})
	c.SetRequestTimeout(time.Duration(TIMEOUT) * time.Second)
	c.Limit(&colly.LimitRule{
		Parallelism: CONCURRENT_URLS,
		DomainGlob:  "*",
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnResponse(func(r *colly.Response) {
		crawlProgress++
		/*
			if crawlProgress > MAXCRAWLVAL {
				log.Println("Stopping crawling due to max URLs visit exceeded!")
				return
			}
		*/
		thisTime := time.Now()
		fmt.Printf("\r%d/%02d/%02d %02d:%02d:%02d Total processed: %d  | Current: %s",
			thisTime.Year(), thisTime.Month(), thisTime.Day(), thisTime.Hour(),
			thisTime.Minute(), thisTime.Second(), crawlProgress, r.Request.URL.String())

		wg := new(sync.WaitGroup)
		// start finding secrets on the html sources
		findSecrets(mainurl, r.Request.URL.String(), &r.Body)
		// start finding secrets within JS files
		wg.Add(1)
		go getJavascript(*r.Request.URL, &r.Body, wg)
		if hastechs {
			executeLoot(p.Host, r.Request.URL.String(), r.StatusCode, &r.Body)
			if SUBMIT_FORM {
				buff := bytes.NewReader(r.Body)
				root, err := html.Parse(buff)
				if err != nil {
					log.Println("error parsing html body:", err.Error(), r.Request.URL.String())
					return
				}
				forms := parseForms(root)
				for _, form := range forms {
					actionURL, err := url.Parse(form.Action)
					if err != nil {
						log.Println(err.Error())
						continue
					}
					actionURL = r.Request.URL.ResolveReference(actionURL)
					mval := setValues(&form.Values)
					err = c.Post(actionURL.String(), mval)
					if err != nil {
						log.Printf("Error posting form %s: %s", r.Request.URL.String(), err.Error())
					}
				}
			}
		}
		wg.Wait()
	})

	c.Visit(mainurl)
	// c.Wait()
	fmt.Print("\n")
}
