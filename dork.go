// 	http://www.google.es/search?start=100&num=100&filter=0&q=site:[QUERY]&gws_rd=cr
// Regex pdf: "https?://.*?\.pdf"

package main

import (
	//	"bytes"
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Received: {"url":"github.com","filetype":["pdf","doc"]}
type Dork struct{
	Url string
	Filetype []string
}


func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	fmt.Println(string(stdin))
	if err != nil {
		log.Println("[/!\\]Error in stdin:", err)
	}

	var dork Dork
	err = json.Unmarshal(stdin,&dork)
	if err != nil {
		log.Println("[/!\\]Error in unmarshal:", err)
	}
	fmt.Println(dork)

	performRequest(dork)
}

func performRequest(dork Dork){
	// Handle number of files
	var filetypes string
	if (len(dork.Filetype)>1){
		filetypes = dork.Filetype[0]
	} else {
		filetypes = dork.Filetype[0]
	}
	// Set the full url
	url := "http://www.google.es/search?start=100&num=100&filter=0&q=site:"+dork.Url+"+filetype:"+filetypes+"&gws_rd=cr"

	// Setup the request to the target
	req, err := http.NewRequest("GET", url, nil)
	// Add the custom headers
	req.Header.Set("User-Agent", "ChromeCaffeine 01.00")

	client := &http.Client{}
	// Perform request and store response on "resp"
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[/!\\]Error in response:", err)
	}   
	defer resp.Body.Close()

	// Handle response (resp)

	// Store HTTP Status code
	fmt.Println(resp.StatusCode)
	// Store Headers
	fmt.Println(resp.Header)
	// Store Body
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

