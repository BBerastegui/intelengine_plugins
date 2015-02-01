
package main

import (
	"encoding/json"
	//	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var(
	initialPath = "/tmp/"
)

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println("[/!\\]Error in stdin:", err)
	}

	type FileList struct{
		Files []string
	}
	var fileList FileList
	err = json.Unmarshal(stdin, &fileList)
	if err != nil {
		log.Println("[/!\\]Error in unmarshal:", err)
	}
	path := initialPath+"test/"
	for _, file := range fileList.Files{
		fmt.Println("[i] Downloading...")
		fmt.Println(file)
		err = download(file, path)
	}

	if err != nil {
		log.Fatalln(err)
	}
}

func download(file string, path string) (err error) {
	log.Println("[i] In performRequest...")

	re := regexp.MustCompile(`.*\/(.*\..*)$`)
	filename := re.FindStringSubmatch(file)
	fmt.Println(path+filename[1])
	out, err := os.Create(path+filename[1])
	defer out.Close()

	req, err := http.NewRequest("GET", file, nil)
	// Add the custom headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.91 Safari/537.36")
	client := &http.Client{}
	// Perform request and store response on "resp"
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Store Body

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
