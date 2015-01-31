// 	http://www.google.es/search?start=100&num=100&filter=0&q=site:[QUERY]&gws_rd=cr
// Regex pdf: "https?://.*?\.pdf"

package main

import (
	//	"bytes"
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"log"
	//	"net/http"
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
}
