package fetchall

import (
	"net/http"
	"os"
	"testing"
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
		t.Log(e.Response.Status)
	}
}
