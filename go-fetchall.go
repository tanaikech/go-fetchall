// Package fetchall (go-fetchall.go) :
// This is a Golang library for running HTTP requests with the asynchronous process.
package fetchall

import (
	"net/http"
	"runtime"
	"sort"
	"sync"
)

const (
	workerNumber = 5 // Default workers
)

// Params : Parameters for fetchAll.
type Params struct {
	Count              int       // Counter of progression.
	Progress           chan int  // Index of request slice running currently.
	Requests           []Request // Requests
	TurnOffMultithread bool      // When this is true, only single thread is used.
	Workers            int       // Default is 5. When 0 is set, all requests are run, simultaneously.
}

// Request : Single request.
type Request struct {
	Client  *http.Client
	Request *http.Request
	index   int
}

// Response : Result
type Response struct {
	Error    error
	Response *http.Response
	index    int
}

// work : For asynchronous process.
type work struct {
	mutex   sync.Mutex
	params  *Params
	result  []Response
	submit  chan Request
	wg      sync.WaitGroup
	workers int
}

// sortResponse : Sort response of slice.
func (w *work) sortResponse() {
	sort.Slice(w.result, func(i, j int) bool { return w.result[i].index < w.result[j].index })
}

// doRequest : Each request is run at this function.
func (w *work) doRequest(r Request) {
	res, err := r.Client.Do(r.Request)
	w.mutex.Lock()
	value := &Response{
		Response: res,
		Error:    err,
		index:    r.index,
	}
	w.result = append(w.result, *value)
	w.mutex.Unlock()
}

// requestHandler : Run requests with asynchronous process.
func (w *work) requestHandler() *work {
	w.submit = make(chan Request, w.workers)
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for {
				p, done := <-w.submit
				if !done {
					return
				}
				w.doRequest(p)
				if w.params.Progress != nil {
					w.params.Progress <- p.index
				}
			}
		}()
	}
	for _, e := range w.params.Requests {
		w.submit <- e
	}
	close(w.submit)
	w.wg.Wait()
	w.sortResponse()
	return w
}

// Do : Executing method
func Do(p *Params) []Response {
	if !p.TurnOffMultithread {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	w := &work{
		workers: func(wn int) int {
			if wn > 0 {
				return wn
			}
			return workerNumber
		}(p.Workers),
		params: p,
	}
	for i := range w.params.Requests {
		w.params.Requests[i].index = i
	}
	res := w.requestHandler()

	// Output : Response from all requests
	return res.result
}
