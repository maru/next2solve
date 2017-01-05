// Next problem to solve
// https://github.com/maru/next2solve
//
// Tests for handlers.go functionality
//
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"next2solve/problems"
	test "next2solve/testing"
	"os"
	"testing"
)

// Valid userid and username values for testing
const (
	userid   = "46232"
	username = "chicapi"
)

var (
	realTest bool
)

// HTTP API test server that responds all requests with an invalid response.
// Wrap for test.InitAPITestServerInvalid function
func initAPITestServerInvalid(t *testing.T, response []string) (*httptest.Server, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(RequestHandler))
	api := test.InitAPITestServerInvalid(t, response)
	problems.InitAPIServer(api.URL)
	return ts, api
}

// HTTP API test server, real API responses were cached in files.
// Wrap for test.InitAPITestServer function
func initAPITestServer(t *testing.T) (*httptest.Server, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(RequestHandler))
	// Test against the real uHunt API web server
	if realTest {
		APIUrl := "http://uhunt.felix-halim.net"
		problems.InitAPIServer(APIUrl)
		return ts, nil
	}
	api := test.InitAPITestServer(t)
	problems.InitAPIServer(api.URL)
	return ts, api
}

// Get the index page
func TestDefaultIndex(t *testing.T) {
	ts, api := initAPITestServer(t)
	defer test.CloseServer(ts)
	defer test.CloseServer(api)

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
	ts, api := initAPITestServer(t)
	defer test.CloseServer(ts)
	defer test.CloseServer(api)

	invalidUsername := "not_" + username
	resp, err := http.PostForm(ts.URL, url.Values{"username": {invalidUsername}, "show-problems": {"ok"}})
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Expected status OK")
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
	ts, api := initAPITestServer(t)
	defer test.CloseServer(ts)
	defer test.CloseServer(api)

	resp, err := http.PostForm(ts.URL, url.Values{"username": {username}, "show-problems": {"ok"}})
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Expected status OK")
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	emtpyError := []byte("<div class=\"error\"></div>")
	if bytes.Index(body, emtpyError) < 0 {
		println(string(body))
		t.Fatal("Expected error empty")
	}
	validUserID := "<input type=\"hidden\" name=\"userid\" value=\"" + userid + "\""
	if bytes.Index(body, []byte(validUserID)) < 0 {
		t.Fatal("Expected userid", userid, "in input text")
	}
}

// Get random problem to solve
func TestRandomProblem(t *testing.T) {
	ts, api := initAPITestServer(t)
	defer test.CloseServer(ts)
	defer test.CloseServer(api)

	resp, err := http.PostForm(ts.URL, url.Values{"username": {username}, "feeling-lucky": {"ok"}})
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Expected status OK")
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

// Get random problem to solve, but nothing to solve
func TestRandomProblemEmpty(t *testing.T) {
	ts, api := initAPITestServerInvalid(t, []string{"[]", userid, "{}"})
	defer test.CloseServer(ts)
	defer test.CloseServer(api)

	resp, err := http.PostForm(ts.URL, url.Values{"username": {username}, "feeling-lucky": {"ok"}})
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Expected status OK")
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	if bytes.Index(body, []byte("<div class=\"error\">No problems to solve</div>")) < 0 {
		println(string(body))
		t.Fatalf("Expected no problems")
	}
}

// Show problems to solve
func TestShowProblemsOk(t *testing.T) {
	ts, api := initAPITestServer(t)
	defer test.CloseServer(ts)
	defer test.CloseServer(api)

	resp, err := http.PostForm(ts.URL, url.Values{"username": {username},
		"show-problems": {"Show problems"}})

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Expected status OK")
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
	if bytes.Index(body, []byte("1337 problems to go!</h2>")) < 0 {
		t.Fatal("Expected problems number")
	}
}

// Test show problems, but nothing to solve
func TestShowProblemsEmpty(t *testing.T) {
	ts, api := initAPITestServerInvalid(t, []string{"[]", userid, "{}"})
	defer test.CloseServer(ts)
	defer test.CloseServer(api)

	resp, err := http.PostForm(ts.URL, url.Values{"username": {username},
		"show-problems": {"Show problems"}})

	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	if bytes.Index(body, []byte("<div class=\"error\">No problems to solve</div>")) < 0 {
		t.Fatalf("Expected no problems")
	}
}

// Initialize the test environment
func TestMain(m *testing.M) {
	flag.BoolVar(&realTest, "real", false, "Test with real uHunt API server")
	flag.Parse()
	os.Exit(m.Run())
}
