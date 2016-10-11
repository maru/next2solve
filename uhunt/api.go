// Next problem to solve
// https://github.com/maru/next2solve
//
// uHunt API calls

package uhunt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// URL paths of uHunt API
const (
	UrlUsernameToUserid  = "/api/uname2uid/%s"
	UrlUserSubmissions   = "/api/subs-user/%s"
	UrlProblemList       = "/api/p"
	UrlProblemListCPBook = "/api/cpbook/%d"
	UrlProblemInfoByID   = "/api/p/id/%d"
	UrlProblemInfoByNum  = "/api/p/num/%d"
)

// Initialize API server with the host URL
func (api *APIServer) Init(url string) {
	api.urlServer = url
}

// Return the host URL
func (api *APIServer) GetUrl() string {
	return api.urlServer
}

// Send request to API and return response body
func (api *APIServer) getResponse(url string) ([]byte, error) {
	resp, err := http.Get(api.urlServer + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Get userid by username, output error if username is not found.
// Returns the userid as a string, or an empty string if error.
func (api *APIServer) GetUserID(username string) (string, error) {
	url := fmt.Sprintf(UrlUsernameToUserid, username)
	resp, err := api.getResponse(url)
	if err != nil {
		return "", err
	}
	id := string(resp)
	return id, nil
}

// Get the problem list of UVa online judge.
// Implemented now: only problems from the CP book 3rd edition.
func (api *APIServer) GetProblemList() ([]APICPBookChapter, error) {
	return api.GetProblemListCPbook(3)
}

// Get the problem list from the CP book from edition 1, 2, or 3.
func (api *APIServer) GetProblemListCPbook(edition int) ([]APICPBookChapter, error) {
	url := fmt.Sprintf(UrlProblemListCPBook, edition)
	resp, err := api.getResponse(url)
	if err != nil {
		return []APICPBookChapter{}, err
	}
	// Parse the data
	var cpProblems []APICPBookChapter
	if err := json.Unmarshal(resp, &cpProblems); err != nil {
		return []APICPBookChapter{}, err
	}

	return cpProblems, nil
}

// Get user submissions
func (api *APIServer) GetUserSubmissions(userid string) (APIUserSubmissions, error) {
	url := fmt.Sprintf(UrlUserSubmissions, userid)
	resp, err := api.getResponse(url)
	if err != nil {
		return APIUserSubmissions{}, err
	}
	var userSubs APIUserSubmissions
	if err := json.Unmarshal(resp, &userSubs); err != nil {
		return APIUserSubmissions{}, err
	}
	userSubs.Submissions = make([]APISubmission, len(userSubs.TmpSubs))
	for i, s := range userSubs.TmpSubs {
		submission := APISubmission{s[0], s[1], s[2], s[3], s[4], s[5], s[6]}
		userSubs.Submissions[i] = submission
	}
	userSubs.TmpSubs = nil
	return userSubs, nil
}

// Get problem information by number
func (api *APIServer) GetProblemByNum(pNumber int) (APIProblem, error) {
	url := fmt.Sprintf(UrlProblemInfoByNum, pNumber)
	resp, err := api.getResponse(url)
	if err != nil {
		return APIProblem{}, err
	}
	var problem APIProblem
	if err := json.Unmarshal(resp, &problem); err != nil {
		return APIProblem{}, err
	}
	return problem, nil
}

// Get problem information by ID
func (api *APIServer) GetProblemByID(pID int) (APIProblem, error) {
	url := fmt.Sprintf(UrlProblemInfoByID, pID)
	resp, err := api.getResponse(url)
	if err != nil {
		return APIProblem{}, err
	}
	var problem APIProblem
	if err := json.Unmarshal(resp, &problem); err != nil {
		return APIProblem{}, err
	}
	return problem, nil
}
