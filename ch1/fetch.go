package main

//fetch prints the content found at a url

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func main() {
	fmt.Println("running")
	for _, url := range os.Args[1:] {
		//add https:// if missing
		prefix, _ := regexp.Match("^https?://", []byte(url))
		if !prefix {
			url = "https://" + url
		}
		//get url
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch:%s\n", err)
			os.Exit(1)

		}
		fmt.Printf("http status:%s\n", resp.Status)
		// b, err := ioutil.ReadAll(resp.Body)
		bytes, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch (%s): error reading body:%s\n", url, err)
			os.Exit(1)
		}
		resp.Body.Close()
		fmt.Printf("bytes copied:%d\n", bytes)
		// fmt.Printf("%s", b)
	}
}
