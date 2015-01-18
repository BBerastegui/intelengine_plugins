package main

import (
	"os"
	"os/exec"
	//	"io"
	"io/ioutil"
	"fmt"
	"encoding/json"
	//	"bytes"
	"net"
	"regexp"
)

type Ip struct{
	Ip net.IP
}

func main() {

	stdin, err := ioutil.ReadAll(os.Stdin)

	var ip Ip
	err = json.Unmarshal(stdin,&ip)
	if err != nil {
		fmt.Println("Error in Unmarshaling...",err)
	}

	fmt.Println(ip.Ip)
	getIpInfo(ip)
}

func getIpInfo(ip Ip){

	cmd := exec.Command("whois", ip.Ip.String())
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	parseWhois(out)
	// First extract the data
	/*
	var w io.Writer
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&out)
	*/
	//fmt.Println(string(out))
}

func parseWhois(whois []byte){
	r, _ := regexp.Compile("^[^%].*:")
	fmt.Println(r.FindString(string(whois)))
}
