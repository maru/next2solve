// Next problem to solve
// https://github.com/maru/next2solve
//
// Tests for problems.go functionality
//
package main

import (
	"next2solve/uhunt"
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
)

// Create the web server and a mock API webserver (reponse is configurable).
func initServer(t *testing.T) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(RequestHandler))
	uhunt.InitAPITestServer(t)
	return ts
}

// Get the index page
func TestDefaultIndex(t *testing.T) {
	ts := initServer(t)
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
	ts := initServer(t)
	defer ts.Close()

	invalidUsername := "not_" + username
	resp, err := http.PostForm(ts.URL, url.Values{"username": {invalidUsername}})
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
	inputUsername := "title=\"Username\" type=\"text\" value=\"" + invalidUsername
	if bytes.Index(body, []byte(inputUsername)) < 0 {
		t.Fatal("Expected username ", invalidUsername, " in input text")
	}
}

// Post a valid username
func TestValidUser(t *testing.T) {
	ts := initServer(t)
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
		t.Fatal("Expected error empty")
	}
	validUserID := "<input type=\"hidden\" name=\"userid\" value=\"" + userid + "\""
	if bytes.Index(body, []byte(validUserID)) < 0 {
		t.Fatal("Expected userid", userid, "in input text")
	}
}

// Check if the userid and username cookies are set
func TestSetCookies(t *testing.T) {
	ts := initServer(t)
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
	ts := initServer(t)
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
		t.Fatal("Expected lucky")
	}

	if bytes.Index(body, []byte("Error template")) >= 0 {
		t.Fatal("Unexpected error")
	}
}

// Get random problem to solve
func TestProblems(t *testing.T) {
	ts := initServer(t)
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
	if bytes.Index(body, []byte("problems")) < 0 {
		t.Fatal("Expected problems")
	}

	if bytes.Index(body, []byte("Error template")) >= 0 {
		t.Fatal("Unexpected error")
	}
}
