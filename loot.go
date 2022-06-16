package main

import (
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

func findSecrets(origPath, secPath string, content *[]byte) {
	if len(regexData) < 1 {
		if !serializeRegexDB(REGEXFILE) {
			log.Fatalln("Error serializing regex data. Exiting...")
		}
	}
	for typerex, regexstr := range regexData {
		demoRex := regexp.MustCompile(regexstr)
		x := demoRex.FindAllSubmatch(*content, -1)
		if len(x) > 0 {
			for _, y := range x {
				// fmt.Print("\n")
				// log.Printf("Secrets found -> Type: %s | Secret: %s", typerex, string(y[0]))
				writer.Write([]string{typerex + " Exposed", origPath, secPath, string(y[0])})
			}
		}
	}
}

func getJavascript(path url.URL, mbody *[]byte, wg *sync.WaitGroup) {
	defer wg.Done()
	xscript := reJSScript.FindAllSubmatch(*mbody, -1)
	if len(xscript) < 1 {
		return
	}
	for _, script := range xscript {
		jsscript := string(script[1])
		murl, err := url.Parse(jsscript)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if !strings.Contains(jsscript, "://") {
			murl = path.ResolveReference(murl)
		}
		_, body, err := makeRequest(murl.String())
		if err != nil {
			log.Println(err.Error())
			continue
		}
		findSecrets(path.String(), murl.String(), body)
	}
}

func executeLoot(domainstacks, path string, status int, body *[]byte) {
	stacks := strings.Split(strings.Split(domainstacks, ":::")[1], ":")
	validated := false
	for _, stack := range stacks {
		for dbstack, sigs := range dbData {
			if dbstack != stack {
				continue
			}
			// initial validation of error/stack trace
			for _, krex := range sigs.Validators.Regex {
				validatorRex := regexp.MustCompile(krex)
				if validateStatusCode(status, sigs.Validators.Status) &&
					validatorRex.Match(*body) {
					validated = true
				}
			}
			if validated {
				for _, kext := range sigs.Extractors {
					mstr := ""
					extractRex := regexp.MustCompile(kext.Regex)
					jkrex := extractRex.FindAllStringSubmatch(string(*body), -1)
					if jkrex != nil {
						for _, kl := range jkrex {
							mstr += strings.Join(kl[1:], " : ") + "\\n"
						}
						writer.Write([]string{kext.Cgroups + " Exposed", path, path, strings.TrimSpace(mstr), strings.Join(stacks, ":")})
					}
				}
			}
		}
	}
}
