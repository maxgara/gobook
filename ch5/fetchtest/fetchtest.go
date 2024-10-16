package fetchtest

//fetch prints the content found at a url

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// Fetch retrieves (HTTP GETs) the resource at url and returns the
// resulting HTTP response body as a string
func Fetch(url string) (string, error) {
	//add https:// if missing
	prefix, _ := regexp.Match("^https?://", []byte(url))
	if !prefix {
		url = "https://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetch: Get %s:%s\n", url, err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return "", fmt.Errorf("fetch: Read body %s:%s\n", url, err)
	}
	resp.Body.Close()
	return string(b), nil
}
