package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func serializeDB(fname string) error {
	dbdata, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dbdata, &dbData); err != nil {
		return err
	}
	return nil
}

func validateStatusCode(code int, codes []int) bool {
	for _, x := range codes {
		if x == code {
			return true
		}
	}
	return false
}

func serializeRegexDB(fname string) bool {
	dbdata, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if err := json.Unmarshal(dbdata, &regexData); err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func makeRequest(url string) (string, *[]byte, error) {
	hostheader := strings.Split(url, "://")[1]
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !verifySSL,
				ServerName:         hostheader,
			},
		},
		Timeout: time.Duration(httpTimeout) * time.Second,
		/*
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		*/
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Host = hostheader
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36`)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`)
	req.Header.Add("Accept-Language", `en-GB,en-US;q=0.9,en;q=0.8`)
	req.Header.Add("Accept-Encoding", `identity`)
	req.Header.Add("DNT", `1`)
	conn, err := client.Do(req)
	if err != nil {
		// log.Printf("Error making request to %s: %s", url, err.Error())
		return "", nil, err
	}
	var respHeaders string
	for key, val := range conn.Header {
		respHeaders += fmt.Sprintf("%s: %s", key, strings.Join(val, ""))
	}
	body, err := ioutil.ReadAll(conn.Body)
	if err != nil {
		return "", nil, err
	}
	return respHeaders, &body, nil
}
