// Next problem to solve
// https://github.com/maru/next2solve
//
// Tests for problems.go functionality
//
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Valid userid and username values for testing
const (
	userid   = "46232"
	username = "chicapi"
)

var (
	idx int
)

// Create the web server and a mock API webserver (reponse is configurable).
func InitServer(apiResponse []string) *httptest.Server {
	idx = 0
	ts := httptest.NewServer(http.HandlerFunc(RequestHandler))
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if idx >= len(apiResponse) {
			panic("Not enough API responses")
		}
		fmt.Fprint(w, apiResponse[idx])
		idx++
	}))
	APIUrl = api.URL
	return ts
}

// Get the index page
func TestDefaultIndex(t *testing.T) {
	ts := InitServer([]string{""})
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	emtpyError := []byte("<div class=\"error\"></div>")
	if bytes.Index(body, emtpyError) < 0 {
		t.Fatal("Expected error empty")
	}
	emptyUsername := []byte("title=\"Username\" type=\"text\" value=\"\"")
	if bytes.Index(body, emptyUsername) < 0 {
		t.Fatal("Expected username empty")
	}
}

// Post an invalid username
func TestInvalidUsername(t *testing.T) {
	ts := InitServer([]string{"0"})
	defer ts.Close()

	username := "not_chicapi"
	resp, err := http.PostForm(ts.URL, url.Values{"username": {username}})
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	notFoundError := []byte("<div class=\"error\">Username not found</div>")
	if bytes.Index(body, notFoundError) < 0 {
		t.Fatal("Expected error 'Username not found'")
	}
	invalidUsername := "title=\"Username\" type=\"text\" value=\"" + username
	if bytes.Index(body, []byte(invalidUsername)) < 0 {
		t.Fatal("Expected username ", username, " in input text")
	}
}

// Post a valid username
func TestValidUser(t *testing.T) {
	ts := InitServer([]string{userid})
	defer ts.Close()

	resp, err := http.PostForm(ts.URL, url.Values{"username": {username}})
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	emtpyError := []byte("<div class=\"error\"></div>")
	if bytes.Index(body, emtpyError) < 0 {
		t.Fatal("Expected error empty", string(body))
	}
	validUserID := "<input type=\"hidden\" name=\"userid\" value=\"" + userid + "\""
	if bytes.Index(body, []byte(validUserID)) < 0 {
		t.Fatal("Expected userid", userid, "in input text")
	}
}

// Check if the userid and username cookies are set
func TestSetCookies(t *testing.T) {
	ts := InitServer([]string{userid})
	defer ts.Close()

	resp, err := http.PostForm(ts.URL, url.Values{"username": {username}})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range resp.Cookies() {
		if c.Name == "userid" && c.Value != userid {
			t.Fatal("Cookie userid value is not", userid, "(", c.Value, ")")
		}
		if c.Name == "username" && c.Value != username {
			t.Fatal("Cookie username value is not", username, "(", c.Value, ")")
		}
	}
	resp, err = http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range resp.Cookies() {
		if c.Name == "userid" && c.Value != userid {
			t.Fatal("Cookie userid value is not", userid, "(", c.Value, ")")
		}
		if c.Name == "username" && c.Value != username {
			t.Fatal("Cookie username value is not", username, "(", c.Value, ")")
		}
	}
}

// Get random problem to solve
func TestRandomProblem(t *testing.T) {
	ts := InitServer([]string{userid})
	defer ts.Close()

	resp, err := http.PostForm(ts.URL, url.Values{"username": {username}, "feeling-lucky": {""}})
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	lucky := []byte("lucky rainbow")
	if bytes.Index(body, lucky) < 0 {
		t.Fatal("Expected lucky", string(body))
	}

	if bytes.Index(body, []byte("Error template")) >= 0 {
		t.Fatal("Unexpected error", string(body))
	}
}

// Get random problem to solve
func TestProblems(t *testing.T) {
	ts := InitServer([]string{userid})
	defer ts.Close()

	resp, err := http.PostForm(ts.URL, url.Values{"username": {username}, "show-problems": {""}})
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	problems := []byte("problems")
	if bytes.Index(body, problems) < 0 {
		t.Fatal("Expected problems", string(body))
	}

	if bytes.Index(body, []byte("Error template")) >= 0 {
		t.Fatal("Unexpected error", string(body))
	}
}
