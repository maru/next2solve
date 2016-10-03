// Next problem to solve
// https://github.com/maru/next2solve
//
// uHunt API calls

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
  "encoding/json"
)

type CPBookChapter struct {
  Title string `json:"title"`
  Subchapters string `json:"arr"`
}

// URL paths of uHunt API
// First argument is the server address, e.g. "http://uhunt.felix-halim.net"
const (
	APIUsernameToUserid          = "%s/api/uname2uid/%s"
	APIUserSubmissions           = "%s/api/subs-user/%s"
	APIUserSubmissionsToProblems = "%s/api/subs-pids/%s/%s/0"
	APIProblemList               = "%s/api/p"
  APIProblemListCPBook         = "%s/api/cpbook/%s"
)

//  Get userid by username, output error if username is not found.
func apiGetUserID(w http.ResponseWriter, username string) (string, error) {
	url := fmt.Sprintf(APIUsernameToUserid, APIUrl, username)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return "", err
	}
	id := string(body)
	if id == "0" {
		return "", errors.New("Username not found")
	}
	return id, nil
}

// Get problem list
func apiGetProblemList() {
  // Implemented now: only problems from the CP book 3rd edition
  return apiGetProblemListCPbook(3)
}

// Get problem list
func apiGetProblemListCPbook(version int) []string {
  var problems []string
  url := fmt.Sprintf(APIProblemListCPBook, APIUrl, version)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Fprintf(w, "Error %v", err)
	// 	return "", err
	// }
  var cp []CPBookChapter
  if err := json.Unmarshal(p, &cp); err != nil {
      return err
  }
	return problems
}

func apiGetUserProblems(userid string) []string {
	var problems []string
	return problems
}
