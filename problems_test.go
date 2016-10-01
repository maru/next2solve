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

// Create mock API webserver. Body is configurable.
func newMockAPIServer(body string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	return ts
}

func TestDefaultIndex(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(RequestHandler))
	defer ts.Close()

	api := newMockAPIServer("")
	defer api.Close()
	APIUrl = api.URL

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	emtpyError := []byte("<div class=\"error\"></div>")
	if i := bytes.Index(body, emtpyError); i < 0 {
		t.Fatal("Expected error empty")
	}
	emptyUsername := []byte("title=\"Username\" type=\"text\" value=\"\"")
	if i := bytes.Index(body, emptyUsername); i < 0 {
		t.Fatal("Expected username empty")
	}
}

func TestInvalidUsername(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(RequestHandler))
	defer ts.Close()

	api := newMockAPIServer("0")
	defer api.Close()
	APIUrl = api.URL

	username := "not_felix_halim"
	res, err := http.PostForm(ts.URL, url.Values{"username": {username}})
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	notFoundError := []byte("<div class=\"error\">Username not found</div>")
	if i := bytes.Index(body, notFoundError); i < 0 {
		t.Fatal("Expected error 'Username not found'")
	}
	invalidUsername := "title=\"Username\" type=\"text\" value=\"" + username
	if i := bytes.Index(body, []byte(invalidUsername)); i < 0 {
		t.Fatal("Expected username ", username, " in input text")
	}
}

func TestValidUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(RequestHandler))
	defer ts.Close()

	userid := "339"
	api := newMockAPIServer(userid)
	defer api.Close()
	APIUrl = api.URL

	username := "felix_halim"
	res, err := http.PostForm(ts.URL, url.Values{"username": {username}})
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	emtpyError := []byte("<div class=\"error\"></div>")
	if i := bytes.Index(body, emtpyError); i < 0 {
		t.Fatal("Expected error empty", string(body))
	}
	validUserID := "<input type=\"hidden\" name=\"userid\" value=\"" + userid + "\""
	if i := bytes.Index(body, []byte(validUserID)); i < 0 {
		t.Fatal("Expected userid", userid, "in input text")
	}
}
