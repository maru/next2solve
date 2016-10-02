/*
http://uhunt.felix-halim.net/api/subs-user/46232

{
  "name": "Maru Berezin",
  "uname": "chicapi",
  "subs":
  [
    [2059528,126,10,1320,1068155560,3,-1],
    [1817422,208,10,5459,1061500105,3,-1]

Submission ID
Problem ID
Verdict ID: 90 : Accepted
Runtime
Submission Time (unix timestamp)
Language ID (1=ANSI C, 2=Java, 3=C++, 4=Pascal, 5=C++11)
Submission Rank
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type UserSubmissions struct {
	Name        string    `json:"name"`
	Username    string    `json:"uname"`
	Submissions [][]int64 `json:"subs"`
}

type Submission struct {
	ID             int64
	ProblemID      int64
	VerdictID      int64
	Runtime        int64
	SubmissionTime time.Time
	LanguageID     int64
	SubmissionRank int64
}

func readUrl(userid string) (body []byte) {
	url := "http://uhunt.felix-halim.net/api/subs-user/" + userid
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func progress() {

	var data []byte
	if len(os.Args) > 1 {
		// Get JSON submissions
		// userid := 46232
		userid := os.Args[1]
		data = readUrl(userid)
	} else {
		var err error
		data, err = ioutil.ReadFile("submissions.json")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Unmarshal
	var userSubs UserSubmissions
	if err := json.Unmarshal(data, &userSubs); err != nil {
		panic(err)
	}

	n := len(userSubs.Submissions)
	// Create statistics
	// problems := make([]int, n)
	// var acProblems []int
	for i := 0; i < n; i++ {
		var s Submission
		s.VerdictID = userSubs.Submissions[i][2]
		// problems[i] =
		s.SubmissionTime = time.Unix(userSubs.Submissions[i][4], 0)
		if s.VerdictID == 90 {
			fmt.Printf("%s\n", s.SubmissionTime)
		}
	}

	// problems by month or week?
}
