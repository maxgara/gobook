package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at`
	Body      string    //in markdown
}
type User struct {
	Login   string
	HTMLURL string `json:"html_url`
}

func main() {
	terms := os.Args[1:]
	results, err := SearchIssues(terms)
	if err != nil {
		fmt.Println("err calling searchIssues: " + err.Error())
		os.Exit(1)
	}
	report := template.Must(template.New("report").
		Funcs(template.FuncMap{"daysAgo": daysAgo}).
		Parse(templ))

	if err := report.Execute(os.Stdout, results); err != nil {
		fmt.Fprintf(os.Stderr, "bad report")
	}

}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	// fmt.Println(q)
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	// fmt.Println(resp)
	// We must close resp.Body on all execution paths
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("http GET resp not ok: %d", resp.StatusCode)
	}
	// fmt.Println(resp.Body)
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	// fmt.Printf("%#v\n", result)
	resp.Body.Close()
	return &result, nil

}

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

const templ = `<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
<th>#</th>
<th>State</th>
<th>User</th>
<th>Title</th>
</tr>
{{range .Items}}
<tr>
<td><a href='{{.HTMLURL}}'>{{.Number}}</td>
<td>{{.State}}</td>
<td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
<td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>
`

type ezwriter string

func (w ezwriter) Write(b []byte) (int, error) {
	w += ezwriter(b)
	return len(b), nil
}

func quickTempHTML(t string, data any) string {
	var w ezwriter
	template.Must(template.New("quickTemp").Parse(t)).Execute(w, data)
	return string(w)
}
