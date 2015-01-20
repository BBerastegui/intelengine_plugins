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
//	"regexp"
)

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
	fmt.Println(string(whois))
}
