// Next problem to solve
// https://github.com/maru/next2solve
//
// Tests for problems.go functionality
//
package uhunt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

//
const (
	nCPBook3Problems = 1658
	nUserProblems    = 319
	problemID        = 1260
	problemNumber    = 10319
	userid           = "46232"
	username         = "chicapi"
)

var (
	apiServer APIServer
)

func InitAPIServer(response string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, response)
	}))
	apiServer.Init(ts.URL)
	return ts
}

func TestGetCPBookProblems(t *testing.T) {
	p, err := ioutil.ReadFile("../data/cpbook3.json")
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	ts := initApiServer(string(p))
	defer ts.Close()
	problems, err := apiServer.GetProblemListCPbook(3)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if len(problems) != nCPBook3Problems {
		t.Fatalf("Error %v", "number of problems does not match")
	}
}

func TestGetCPBookProblemsReal(t *testing.T) {
	return

	APIUrl := "http://uhunt.felix-halim.net"
	apiServer.Init(APIUrl)
	problems, err := apiServer.GetProblemListCPbook(3)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if len(problems) != nCPBook3Problems {
		t.Fatalf("Error %v", "number of problems does not match")
	}
}

func TestGetUserProblems(t *testing.T) {
	p, err := ioutil.ReadFile("../test/submissions.json")
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	ts := initApiServer(string(p))
	defer ts.Close()
	problems, err := apiServer.GetUserProblems(userid)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if len(problems) != nUserProblems {
		t.Fatalf("Error %v %d", "number of problems does not match", len(problems))
	}
}

func TestGetProblemInfo(t *testing.T) {
	p, err := ioutil.ReadFile("../test/p1260.json")
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	ts := initApiServer(string(p))
	defer ts.Close()
	problem, err := apiServer.GetProblemInfo(problemID)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if problem.ProblemNumber != problemNumber {
		t.Fatalf("Error %v", "problem id does not match")
	}
}
