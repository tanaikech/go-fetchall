package fetchall

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"golang.org/x/net/html"
)

func TestRunFetchAll(t *testing.T) {
	// URLs for requesting
	urls := []string{"https://tanaikech.github.io/"}

	p := &Params{}
	for _, e := range urls {
		req, err := http.NewRequest("GET", e, nil)
		if err != nil {
			os.Exit(1)
		}
		r := &Request{
			Request: req,
			Client:  &http.Client{},
		}
		p.Requests = append(p.Requests, *r)
	}

	// Run fetchall
	res := Do(p)

	// Show result
	for _, e := range res {
		doc, err := html.Parse(e.Response.Body)
		if err != nil {
			os.Exit(1)
		}
		var title *html.Node
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "title" {
				title = n
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
		var bf bytes.Buffer
		html.Render(&bf, title)
		t.Log(bf.String())
	}
}
