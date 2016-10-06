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

// HTTP API test server that responds all requests with an empty string.
func InitAPITestServerInvalid(t *testing.T, apiServer *APIServer) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response string
		switch r.RequestURI {
		default:
			println(r.RequestURI)
			response = ""
		}
		fmt.Fprint(w, response)
	}))
	apiServer.Init(ts.URL)
	return ts
}

// HTTP API test server, real API responses were cached in files.
func InitAPITestServer(t *testing.T, apiServer *APIServer) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response string
		switch r.RequestURI {
		case "/api/uname2uid/not_" + username:
			response = "0"
		case "/api/uname2uid/" + username:
			response = userid
		case "/api/cpbook/3":
			problems, err := ioutil.ReadFile("../data/cpbook3.json")
			if err != nil {
				t.Fatalf("Error %v", err)
			}
			response = string(problems)
		default:
			println(r.RequestURI)
			response = "{}"
			t.Fatalf("Error %v", "err")
		}
		fmt.Fprint(w, response)
	}))
	apiServer.Init(ts.URL)
	return ts
}

// Test API Server URL
func TestGetUrl(t *testing.T) {
	var apiServer APIServer
	ts := InitAPITestServer(t, &apiServer)
	defer ts.Close()

	if ts.URL != apiServer.GetUrl() {
		t.Fatalf("Expected API server URL %s, got %s", ts.URL, apiServer.GetUrl())
	}
}

// Test invalid username
func TestInvalidUsername(t *testing.T) {
	var apiServer APIServer
	ts := InitAPITestServer(t, &apiServer)
	defer ts.Close()

	invalidUsername := "not_" + username
	id, err := apiServer.GetUserID(invalidUsername)

	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if id != "0" {
		t.Fatalf("Error expected userid 0")
	}
}

// Test username, invalid response
func TestUsernameInvalidResponse(t *testing.T) {
	var apiServer APIServer
	ts := InitAPITestServerInvalid(t, &apiServer)
	defer ts.Close()

	id, err := apiServer.GetUserID(username)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if id != "" {
		t.Fatalf("Error expected an empty response")
	}
}

// Test valid username
func TestValidUsername(t *testing.T) {
	var apiServer APIServer
	ts := InitAPITestServer(t, &apiServer)
	defer ts.Close()

	id, err := apiServer.GetUserID(username)

	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if id != userid {
		t.Fatalf("Error expected userid %s", userid)
	}
}

// Test get problems from CP book
func TestGetCPBookProblems(t *testing.T) {
	var apiServer APIServer
	ts := InitAPITestServer(t, &apiServer)
	defer ts.Close()

	problems, err := apiServer.GetProblemListCPbook(3)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if problems[0].Title != "Introduction" {
		t.Fatalf("Error %v", "title does not match")
	}
}

// Test get problems from CP book, invalid response
func TestGetCPBookProblemsInvalidResponse(t *testing.T) {
	var apiServer APIServer
	ts := InitAPITestServerInvalid(t, &apiServer)
	defer ts.Close()

	problems, err := apiServer.GetProblemListCPbook(3)
	if err == nil {
		t.Fatalf("Error %v", "expected end of JSON input")
	}
	if len(problems) != 0 {
		t.Fatalf("Error expected an empty response")
	}
}


//
// func TestGetCPBookProblemsReal(t *testing.T) {
// 	return
//
// 	APIUrl := "http://uhunt.felix-halim.net"
// 	apiServer.Init(APIUrl)
// 	problems, err := apiServer.GetProblemListCPbook(3)
// 	if err != nil {
// 		t.Fatalf("Error %v", err)
// 	}
// 	if len(problems) != nCPBook3Problems {
// 		t.Fatalf("Error %v", "number of problems does not match")
// 	}
// }
//
// func TestGetUserSubmissions(t *testing.T) {
// 	ts := InitAPIServer(t)
// 	defer ts.Close()
// 	problems, err := apiServer.GetUserSubmissions(userid)
// 	if err != nil {
// 		t.Fatalf("Error %v", err)
// 	}
// 	if len(problems) != nUserProblems {
// 		t.Fatalf("Error %v %d", "number of problems does not match", len(problems))
// 	}
// }
//
// func TestGetProblemInfo(t *testing.T) {
// 	ts := InitAPIServer(t)
// 	defer ts.Close()
// 	problem, err := apiServer.GetProblemInfoByNum(problemNumber)
// 	if err != nil {
// 		t.Fatalf("Error %v", err)
// 	}
// 	if problem.ProblemNumber != problemNumber {
// 		t.Fatalf("Error %v", "problem id does not match")
// 	}
// }
