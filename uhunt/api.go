// Next problem to solve
// https://github.com/maru/next2solve
//
// uHunt API calls

package uhunt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

// URL paths of uHunt API
const (
	APIUsernameToUserid  = "/api/uname2uid/%s"
	APIUserSubmissions   = "/api/subs-user/%s"
	APIProblemList       = "/api/p"
	APIProblemListCPBook = "/api/cpbook/%d"
	APIProblemInfoByNum  = "/api/p/num/%d"
)

// Initialize API server with the host URL
func (api *APIServer) Init(url string) {
	api.urlServer = url
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
	url := fmt.Sprintf(APIUsernameToUserid, username)
	resp, err := api.getResponse(url)
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
// Implemented now: only problems from the CP book 3rd edition
func (api *APIServer) GetProblemList() ([]int, error) {
	return api.GetProblemListCPbook(3)
}

// Get problem list
func (api *APIServer) GetProblemListCPbook(version int) ([]int, error) {
	var problems []int
	url := fmt.Sprintf(APIProblemListCPBook, version)
	resp, err := api.getResponse(url)
	if err != nil {
		return problems, err
	}
	// Parse the data
	var cpBook []CPBookChapter
	if err := json.Unmarshal(resp, &cpBook); err != nil {
		return problems, err
	}
	// Create an array with the problem ids
	numChapters := 0
	numSubChapters := 100
	numSubSubChapters := 1000
	for _, chapter := range cpBook {
		numChapters++
		for _, subchapter := range chapter.Subchapters {
			numSubChapters++
			for _, subsubchapter := range subchapter.Subsubchapters {
				numSubSubChapters++
				arr := subsubchapter.([]interface{})
				for _, p := range arr[1:] {
					pid := int(math.Abs(p.(float64)))
					problems = append(problems, pid)
					api.cpProblems[pid] = CPProblem{p.(float64) < 0, []int{numChapters, numSubChapters, numSubSubChapters}}
				}
			}
		}
	}
	return problems, nil
}

// Get user submissions, so we can obtain the solved problems
func (api *APIServer) GetUserProblems(userid string) (map[int]bool, error) {
	url := fmt.Sprintf(APIUserSubmissions, userid)
	resp, err := api.getResponse(url)
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

func (api *APIServer) GetProblemInfoByNum(pnum int) (Problem, error) {
	url := fmt.Sprintf(APIProblemInfoByNum, pnum)
	resp, err := api.getResponse(url)
	if err != nil {
		return Problem{}, err
	}
	var problem Problem
	if err := json.Unmarshal(resp, &problem); err != nil {
		return Problem{}, err
	}
	// ProblemInfo
	return problem, nil
}
