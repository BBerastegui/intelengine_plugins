// http://www.google.es/search?start=100&num=100&filter=0&q=site:[QUERY]&gws_rd=cr
// Regex pdf: "https?://.*?\.pdf"

package main

import (
	//	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var (
	reNumResults = regexp.MustCompile(`<span class="sb_count">.*?([\d.,]+).*?</span>`)
	reNonNumbers = regexp.MustCompile(`\D+`)
)

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println("[/!\\]Error in stdin:", err)
	}

	var dork Dork
	err = json.Unmarshal(stdin, &dork)
	if err != nil {
		log.Println("[/!\\]Error in unmarshal:", err)
	}

	log.Println("[i] Performing request...")
	rp := NewResultsParser(dork)
	//results, err := rp.DoRequest()
	_, err = rp.DoRequest()
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("results=%#v\n", results)
}

// Received: {"url":"github.com","filetype":["pdf","doc"]}
type Dork struct {
	Url      string
	Filetype []string
}

type ResultsParser struct {
	dork Dork
}

func NewResultsParser(dork Dork) *ResultsParser {
	return &ResultsParser{dork: dork}
}

func (rp *ResultsParser) DoRequest() (results []string, err error) {
	log.Println("[i] In performRequest...")
	results = []string{}
	for _, ft := range rp.dork.Filetype {
		partialResults, err := rp.requestByFiletype(ft)
		if  err != nil {
			return nil, err
		}
		//fmt.Printf("partialResults=%#v\n", partialResults)
		results = append(results, partialResults...)
	}
	return results, nil
}

func (rp *ResultsParser) requestByFiletype(filetype string) (results []string, err error) {
	reFileURL, err := regexp.Compile(`href="(https?://[^"]*?\.` + filetype + `)`)
	if err != nil {
		return nil, err
	}
	results = []string{}
	firstPage := true
	numResults, curPage := 0, 1
	for {
		curPageStr := strconv.Itoa(curPage)
		if curPageStr == "" {
			return nil, errors.New("cannot convert page number to string")
		}
		url := "http://www.bing.com/search?first=" + curPageStr + "&q=site:" + rp.dork.Url + "%20filetype:" + filetype

		// Setup the request to the target
		req, err := http.NewRequest("GET", url, nil)
		// Add the custom headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.91 Safari/537.36")
		// Set Cookie
		req.AddCookie(&http.Cookie{Name: "SRCHHPGUSR", Value: "ADLT=OFF&NRSLT=100"})
		req.AddCookie(&http.Cookie{Name: "MUID", Value: "00000000000000000000000000000000"})
		client := &http.Client{}
		// Perform request and store response on "resp"
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		// Store Body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if firstPage {
			sm := reNumResults.FindStringSubmatch(string(body))
			if len(sm) != 2 {
				return nil, errors.New("number of results not found")
			}
			numResults, err = strconv.Atoi(reNonNumbers.ReplaceAllString(sm[1], ""))
			if err != nil {
				return nil, err
			}
			firstPage = false
		}
		pageResults := reFileURL.FindAllStringSubmatch(string(body), -1)
		if len(pageResults) == 0 {
			break
		}
		for _, pr := range pageResults {
			results = append(results, pr[1])
			fmt.Println(pr[1])
		}
		//fmt.Printf("partialFiletypeResults=%s ; %s\n",pageResults[0][1], pageResults[len(pageResults)-1][1])
		log.Printf("page=%d numresults=%d\n", curPage, numResults)

		curPage += 10
		if curPage > numResults || curPage > 1500 {
			break
		}
	}
	return results, nil
}
