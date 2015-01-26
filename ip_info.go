package main

import (
	"os"
	"os/exec"
	//	"io"
	"io/ioutil"
	"fmt"
	"encoding/json"
	//	"bytes"
	//	"net"
	"regexp"
)

type Person struct{
	name string
	surname string
	address string
}

type Domain struct{
	domainName string
	registrar string
	sponsoringRegistrar string
}

func main() {

	stdin, err := ioutil.ReadAll(os.Stdin)

	var dat map[string]interface{}
	err = json.Unmarshal(stdin,&dat)
	if err != nil {
		fmt.Println("Error in Unmarshaling...",err)
	}

	if (dat["domain"] != nil) {
		runWhois(dat["domain"].(string))
	} else {
		runWhois(dat["ip"].(string))
	}
}

func runWhois(dat string){

	cmd := exec.Command("whois", dat)
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	parseWhois(out)
}

func parseWhois(whois []byte){

	var r map[string]*regexp.Regexp
	r = make(map[string]*regexp.Regexp)

	// First extract Domain Admin information
	r["adminName"] = regexp.MustCompile("Admin Name:(.*)")
	r["adminOrganization"] = regexp.MustCompile("Admin Organization:(.*)")
	r["adminStreet"] = regexp.MustCompile("Admin Street:(.*)")
	r["adminCity"] = regexp.MustCompile("Admin City:(.*)")

	/*
	for k, _ := range r {
		fmt.Println(r[k].FindStringSubmatch(string(whois)))
	}
	*/

	jsonString, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(jsonString)

}

