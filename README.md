# go-fetchall

[![Build Status](https://travis-ci.org/tanaikech/go-fetchall.svg?branch=master)](https://travis-ci.org/tanaikech/go-fetchall)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENCE)

<a name="top"></a>

# Overview

This is a Golang library for running HTTP requests with the asynchronous process. The progress of requests can be also shown.

# Demo

![](images/demo.gif)

In this demonstration, 5 requests are run by 2 workers. And before each request, the waiting time for 2 seconds is added as a sample. By this, you can easily see the work with 2 workers. Also you can see this script at [the following sample script](#samplescript).

# Description

Recently, when I created applications, I had seen the situation which was required to run HTML requests with the asynchronous process. When several requests were run, it was required to know the progression of requests. So I created this as a library.

## Features

- This library can run several HTTP requests with the asynchronous process.
- The progression of requests can be retrieved.

# Install

You can install this using `go get` as follows.

```bash
$ go get -u github.com/tanaikech/go-fetchall
```

# Method

| Method             | Explanation                                          |
| :----------------- | :--------------------------------------------------- |
| `Do(*http.Client)` | Run inputted requests with the asynchronous process. |

# Usage

### Set parameters

At first, users are required to set the parameters for using this library. You can see about this at the following sample script.

#### Input value:

```golang
// Params : Parameters for fetchall.
type Params struct {
	Count              int       // Counter of progression.
	Progress           chan int  // Index of request slice running currently.
	Requests           []Request // Requests
	TurnOffMultithread bool      // When this is true, only single thread is used.
	Workers            int       // Default is 5. When 0 is set, all requests are run, simultaneously.
}
```

#### Response value:

A slice with the following structure is returned.

```golang
// Response : Result
type Response struct {
	Error    error
	Response *http.Response
	index    int
}
```

<a name="samplescript"></a>

### Sample script 1

In this sample script, 5 requests are run by 2 workers. The title tag is retrieved from the response. The progress of requests are shown.

```golang
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	fetchall "github.com/tanaikech/go-fetchall"
	"golang.org/x/net/html"
)

func main() {
	// URLs for requesting. There are sample URLs in this case.
	urls := []string{
		"https://tanaikech.github.io/2019/04/30/split-array-by-n-elements-using-google-apps-script/",
		"https://tanaikech.github.io/2019/04/20/overwriting-several-google-documents-by-2-text-files-using-google-apps-script/",
		"https://tanaikech.github.io/2019/04/20/creating-google-document-by-converting-pdf-and-image-files-with-ocr-using-google-apps-script/",
		"https://tanaikech.github.io/2019/04/20/gas-library---fetchapp/",
		"https://tanaikech.github.io/2019/04/11/converting-many-files-to-google-docs-using-google-apps-script/",
	}

	p := &fetchall.Params{
		Count:    len(urls),
		Progress: make(chan int),
		Workers:  2,
  }

  // Create requests
	for _, e := range urls {
		req, err := http.NewRequest("GET", e, nil)
		if err != nil {
			os.Exit(1)
		}
		r := &fetchall.Request{
			Request: req,
			Client:  &http.Client{},
		}
		p.Requests = append(p.Requests, *r)
	}

	// Show progression of requests
	go func() {
		for {
			v, done := <-p.Progress
			if !done {
				return
			}
			p.Count--
			fmt.Printf("index: %d, counter: %d\n", v, p.Count)
		}
	}()

	// Run fetchall
	res := fetchall.Do(p)
	close(p.Progress)

	// Show result. As a sample, the title tags are retrieved from each URL.
	fmt.Printf("\n- Result\n")
	for _, e := range res {

		// Ref: https://stackoverflow.com/a/38855264
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
		fmt.Printf("%s\n", bf.String())
	}
}
```

### Sample script 2

In this sample script, 5 requests are run by 2 workers. The title tag is retrieved from the response. The progress of requests are **NOT** shown.

```golang
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	fetchall "github.com/tanaikech/go-fetchall"
	"golang.org/x/net/html"
)

func main() {
	// URLs for requesting
	urls := []string{
		"https://tanaikech.github.io/2019/04/30/split-array-by-n-elements-using-google-apps-script/",
		"https://tanaikech.github.io/2019/04/20/overwriting-several-google-documents-by-2-text-files-using-google-apps-script/",
		"https://tanaikech.github.io/2019/04/20/creating-google-document-by-converting-pdf-and-image-files-with-ocr-using-google-apps-script/",
		"https://tanaikech.github.io/2019/04/20/gas-library---fetchapp/",
		"https://tanaikech.github.io/2019/04/11/converting-many-files-to-google-docs-using-google-apps-script/",
	}

	p := &fetchall.Params{}
	for _, e := range urls {
		req, err := http.NewRequest("GET", e, nil)
		if err != nil {
			os.Exit(1)
		}
		r := &fetchall.Request{
			Request: req,
			Client:  &http.Client{},
		}
		p.Requests = append(p.Requests, *r)
	}

	// Run fetchall
	res := fetchall.Do(p)

	// Show result
	for _, e := range res {

		// Ref: https://stackoverflow.com/a/38855264
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
		fmt.Printf("%s\n", bf.String())
	}
}
```

---

<a name="Licence"></a>

# Licence

[MIT](LICENCE)

<a name="Author"></a>

# Author

[Tanaike](https://tanaikech.github.io/about/)

If you have any questions and commissions for me, feel free to tell me.

<a name="Update_History"></a>

# Update History

- v1.0.0 (May 1, 2019)

  1. Initial release.

[TOP](#top)
