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

// HTTP API test server, real API responses were cached in files.
// Wrap for test.InitAPITestServer function
func initAPITestServer(t *testing.T) *httptest.Server {
	// Test against the real uHunt API web server
	if realTest {
		APIUrl := "http://uhunt.felix-halim.net"
		InitAPIServer(APIUrl)
		return nil
	}
	api := test.InitAPITestServer(t)
	InitAPIServer(api.URL)
	return api
}

// Test initialize API server
func TestInitAPIServer(t *testing.T) {
	api := initAPITestServer(t)
	defer test.CloseServer(api)

	InitAPIServer(api.URL)
	if api.URL != apiServer.GetUrl() {
		t.Fatalf("Expected API server URL %s, got %s", api.URL, apiServer.GetUrl())
	}
}

// Test get userid with invalid username
func TestGetUserIDInvalid(t *testing.T) {
	api := initAPITestServer(t)
	defer test.CloseServer(api)

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
	api := initAPITestServer(t)
	defer test.CloseServer(api)

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
	api := initAPITestServer(t)
	defer test.CloseServer(api)

	problems := GetUnsolvedProblemsCPBook("0")
	if len(problems) != 0 {
		t.Fatalf("Expected empty problem list")
	}
}

// Initialize the test environment
func TestMain(m *testing.M) {
	flag.BoolVar(&realTest, "real", false, "Test with real uHunt API server")
	flag.Parse()
	os.Exit(m.Run())
}
