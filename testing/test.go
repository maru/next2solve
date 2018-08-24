// Next problem to solve
// https://github.com/maru/next2solve
//
// Helper test functions

package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"testing"
)

var (
	idx int
)

// Close the test web server
func CloseServer(ts *httptest.Server) {
	if ts != nil {
		ts.Close()
	}
}

// HTTP API test server that responds all requests with an invalid response.
func InitAPITestServerInvalid(t *testing.T, response []string) *httptest.Server {
	idx = 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if idx >= len(response) {
			t.Fatal("Not enough responses!")
		}
		fmt.Fprint(w, response[idx])
		idx++
	}))
	return ts
}

// HTTP API test server, real API responses were cached in files.
func InitAPITestServer(t *testing.T) *httptest.Server {
	// Get base directory
	_, filename, _, _ := runtime.Caller(0)
	baseDirectory := path.Dir(filename)
	// Create a test API web server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, err := ioutil.ReadFile(path.Join(baseDirectory, r.RequestURI))
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		fmt.Fprint(w, string(response))
	}))
	return ts
}
