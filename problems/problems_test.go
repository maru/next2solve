// Next problem to solve
// https://github.com/maru/next2solve
//
// Tests for problems.go functionality
//
package problems

import (
	// 	"bytes"
	"flag"
	"net/http/httptest"
	test "next2solve/testing"
	"os"
	"testing"
)

//
const (
// nCPBook3Problems  = 1658
// nUnsolvedProblems = 1589
userid   = "46232"
username = "chicapi"
)

var (
	realTest bool
)

// HTTP API test server that responds all requests with an invalid response.
// Wrap for test.InitAPITestServerInvalid function
func initAPITestServerInvalid(t *testing.T, response []string) *httptest.Server {
	ts := test.InitAPITestServerInvalid(t, response)
	InitAPIServer(ts.URL)
	return ts
}

// HTTP API test server, real API responses were cached in files.
// Wrap for test.InitAPITestServer function
func initAPITestServer(t *testing.T) *httptest.Server {
	// Test against the real uHunt API web server
	if realTest {
		APIUrl := "http://uhunt.felix-halim.net"
		InitAPIServer(APIUrl)
		return nil
	}
	ts := test.InitAPITestServer(t)
	InitAPIServer(ts.URL)
	return ts
}

// Test initialize API server
func TestInitAPIServer(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	InitAPIServer(ts.URL)
	if ts.URL != apiServer.GetUrl() {
		t.Fatalf("Expected API server URL %s, got %s", ts.URL, apiServer.GetUrl())
	}
}

// Test get userid with invalid username
func TestGetUserIDInvalid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	invalidUsername := "not_" + username
	id, err := GetUserID(invalidUsername)
	if err == nil || err.Error() != "Username not found" {
		t.Fatalf("Expected error")
	}
	if id != "" {
		t.Fatalf("Expected empty userid")
	}
}

// Test get userid with valid username
func TestGetUserIDValid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	id, err := GetUserID(username)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	if id != userid {
		t.Fatalf("Expected userid %v got %v", userid, id)
	}
}

// Test GetUnsolvedProblemsCPBook, invalid userid
func TestGetUnsolvedProblemsCPBookInvalidUserid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblemsCPBook("0")
	if len(problems) != 0 {
		t.Fatalf("Expected empty problem list")
	}
}

// Test GetUnsolvedProblemsCPBook, valid userid
func TestGetUnsolvedProblemsCPBookValidUserid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblemsCPBook(userid)
	if len(problems) == 0 {
		t.Fatalf("Expected problem list")
	}
}

// Test GetUnsolvedProblemsCPBook, valid userid
func TestGetUnsolvedProblemsCPBookInvalidResponse(t *testing.T) {
	ts := initAPITestServerInvalid(t, []string{"", ""})
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblemsCPBook(userid)
	if len(problems) != 0 {
		t.Fatalf("Expected empty problem list")
	}
}

// Test GetUnsolvedProblems, valid userid
func TestGetUnsolvedProblemsValidUserid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblems(userid)
	if len(problems) == 0 {
		t.Fatalf("Expected problem list")
	}
}

// Test GetUnsolvedProblemRandom, invalid userid
func TestGetUnsolvedProblemRandomInvalidUserid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblemRandom("0")
	if len(problems) != 0 {
		t.Fatalf("Expected empty problem list")
	}
}

// Test GetUnsolvedProblemRandom, valid userid
func TestGetUnsolvedProblemRandomValidUserid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblemRandom(userid)
	if len(problems) == 0 {
		t.Fatalf("Expected problem list")
	}
}

// Initialize the test environment
func TestMain(m *testing.M) {
	flag.BoolVar(&realTest, "real", false, "Test with real uHunt API server")
	flag.Parse()
	os.Exit(m.Run())
}
