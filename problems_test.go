// Next problem to solve
// https://github.com/maru/next2solve
//
// Tests for problems.go functionality
//
package main

import (
	// 	"bytes"
	// 	"fmt"
	"io/ioutil"
	"testing"
)

//
const (
	nCPBook3Problems = 1658
	nUnsolvedProblems = 1589
)

var ()

// Init API server with problems and submissions from testing files
func loadAPIProblems(t *testing.T) []string {
	problems, err := ioutil.ReadFile("data/cpbook3.json")
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	submissions, err := ioutil.ReadFile("test/submissions.json")
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	responses := []string{string(problems), string(submissions)}
	problemInfo, err := ioutil.ReadFile("test/p1260.json")
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	for i := 0; i < nCPBook3Problems; i++ {
		responses = append(responses, string(problemInfo))
	}
	return responses
}

func TestGetUnsolvedProblems(t *testing.T) {
	resp := loadAPIProblems(t)
	ts := initApiServer(t, resp)
	defer ts.Close()

	unsolved := GetUnsolvedProblems(userid)
	println(len(unsolved))
	if len(unsolved) == 0 {
		t.Fatalf("Error %v", "no problems to solve!")
	}
}
