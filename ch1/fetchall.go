package main

//fetch prints the content found at a url

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	// "regexp"
	"time"
)

func main() {
	ch := make(chan string)

	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}

}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, _ := http.Get(url)
	bytes, _ := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, bytes, url)
	//add https:// if missing
	// prefix, _ := regexp.Match("^https?://", []byte(url))
	// if !prefix {
	// 	url = "https://" + url
	// }
	//get url
	// resp, err := http.Get(url)
	// if err != nil {
	// 	ch <- fmt.Sprint(err)
	// 	return
	// }
	// bytes, err := io.Copy(ioutil.Discard, resp.Body)

	// if err != nil {
	// 	ch <- fmt.Sprintf("error copying %s,%v\n", url, err)

	// }
	// fmt.Printf("%s", b)
}
