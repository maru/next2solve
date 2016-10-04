// Next problem to solve
// https://github.com/maru/next2solve
//
// uHunt API calls

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

// URL paths of uHunt API
// First argument is the server address, e.g. "http://uhunt.felix-halim.net"
const (
	APIUsernameToUserid = "%s/api/uname2uid/%s"
	APIUserSubmissions  = "%s/api/subs-user/%s"
	// APIUserSubmissionsToProblems = "%s/api/subs-pids/%s/%s/0"
	APIProblemList       = "%s/api/p"
	APIProblemListCPBook = "%s/api/cpbook/%d"
	APIProblemInfo       = "%s/api/p/id/%d"
)

// Send request to API and return response body
func getResponse(url string) ([]byte, error) {
	resp, err := http.Get(url)
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

//  Get userid by username, output error if username is not found.
func apiGetUserID(username string) (string, error) {
	url := fmt.Sprintf(APIUsernameToUserid, APIUrl, username)
	resp, err := getResponse(url)
	if err != nil {
		return "", err
	}
	id := string(resp)
	if id == "0" {
		return "", errors.New("Username not found")
	}
	return id, nil
}

// Get problem list
func apiGetProblemList() ([]int, error) {
	// Implemented now: only problems from the CP book 3rd edition
	return apiGetProblemListCPbook(3)
}

type CPBookChapter struct {
	Title       string             `json:"title"`
	Subchapters []CPBookSubchapter `json:"arr"`
}
type CPBookSubchapter struct {
	Title          string        `json:"title"`
	Subsubchapters []interface{} `json:"arr"`
}

// Get problem list
func apiGetProblemListCPbook(version int) ([]int, error) {
	var problems []int
	url := fmt.Sprintf(APIProblemListCPBook, APIUrl, version)
	resp, err := getResponse(url)
	if err != nil {
		return problems, err
	}
	var cpBook []CPBookChapter
	if err := json.Unmarshal(resp, &cpBook); err != nil {
		return problems, err
	}
	for _, chapter := range cpBook {
		for _, subchapter := range chapter.Subchapters {
			for _, subsubchapter := range subchapter.Subsubchapters {
				arr := subsubchapter.([]interface{})
				for _, p := range arr[1:] {
					problems = append(problems, int(p.(float64)))
				}
			}
		}
	}
	return problems, nil
}

type UserSubmissions struct {
	Name        string      `json:"name"`
	Username    string      `json:"uname"`
	Submissions [][]float64 `json:"subs"`
}

// Get user submissions, so we can obtain the solved problems
func apiGetUserProblems(userid string) (map[int]bool, error) {
	url := fmt.Sprintf(APIUserSubmissions, APIUrl, userid)
	resp, err := getResponse(url)
	if err != nil {
		return nil, err
	}
	var userSubs UserSubmissions
	if err := json.Unmarshal(resp, &userSubs); err != nil {
		return nil, err
	}

	// Get only accepted (distinct) problems
	// An element in the array contains:
	//  0   Submission ID
	//  1   Problem ID  (* we want this)
	//  2   Verdict ID  (* and this with value 90 : Accepted)
	//  3   Runtime
	//  4   Submission Time (unix timestamp)
	//  5   Language ID (1=ANSI C, 2=Java, 3=C++, 4=Pascal, 5=C++11)
	//  6   Submission Rank
	solved := make(map[int]bool)
	for _, p := range userSubs.Submissions {
		pid := int(math.Abs(p[1]))
		if p[2] == 90 {
			solved[pid] = true
		}
	}
	return solved, nil
}

type Problem struct {
	ProblemID              int64  `json:"pid"`
	ProblemNumber          int64  `json:"num"`
	Title                  string `json:"title"`
	Dacu                   int64  `json:"dacu"`
	BestRuntime            int64  `json:"mrun"`
	BestUsedMemory         int64  `json:"mmem"`
	NumNoVerdict           int64  `json:"nover"`
	NumSubmissionError     int64  `json:"sube"`
	NumCantBeJudged        int64  `json:"noj"`
	NumInQueue             int64  `json:"inq"`
	NumCompilationError    int64  `json:"ce"`
	NumRestrictedFunction  int64  `json:"rf"`
	NumRuntimeError        int64  `json:"re"`
	NumOutputLimitExceeded int64  `json:"ole"`
	NumTimeLimitExceeded   int64  `json:"tle"`
	NumMemoryLimitExceeded int64  `json:"mle"`
	NumWrongAnswer         int64  `json:"wa"`
	NumPresentationError   int64  `json:"pe"`
	NumAccepted            int64  `json:"ac"`
	RunTimeLimit           int64  `json:"rtl"`
	Status                 int64  `json:"status"`
	Rej                    int64  `json:"rej"`
}

func apiGetProblemInfo(pid int) (interface{}, error) {
	url := fmt.Sprintf(APIProblemInfo, APIUrl, pid)
	resp, err := getResponse(url)
	if err != nil {
		return nil, err
	}
	var problem Problem
	if err := json.Unmarshal(resp, &problem); err != nil {
		return nil, err
	}
	// ProblemInfo
	return problem, nil
}
