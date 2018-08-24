// Next problem to solve
// https://github.com/maru/next2solve
//
// Tests for API functionality
//
// To test against the real uHunt API web server:
//   go test next2solve/uhunt -args -real

package uhunt

import (
	"flag"
	"net/http/httptest"
	test "next2solve/testing"
	"os"
	"testing"
)

const (
	nUserSubmissions = 729
	problemID        = 1260
	problemNumber    = 10319
	userid           = "46232"
	username         = "chicapi"
)

var (
	realTest bool
)

// HTTP API test server that responds all requests with an invalid response.
// Wrap for test.InitAPITestServerInvalid function
func initAPITestServerInvalid(t *testing.T, apiServer *APIServer, response string) *httptest.Server {
	ts := test.InitAPITestServerInvalid(t, []string{response})
	apiServer.Init(ts.URL)
	return ts
}

// HTTP API test server, real API responses were cached in files.
// Wrap for test.InitAPITestServer function
func initAPITestServer(t *testing.T, apiServer *APIServer) *httptest.Server {
	// Test against the real uHunt API web server
	if realTest {
		APIUrl := "https://uhunt.onlinejudge.org"
		apiServer.Init(APIUrl)
		return nil
	}
	ts := test.InitAPITestServer(t)
	apiServer.Init(ts.URL)
	return ts
}

// Test API Server URL
func TestGetUrl(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServer(t, &apiServer)
	defer test.CloseServer(ts)

	url := "https://uhunt.onlinejudge.org"
	if ts != nil {
		url = ts.URL
	}
	if url != apiServer.GetUrl() {
		t.Fatalf("Expected API server URL %s, got %s", url, apiServer.GetUrl())
	}
}

// Test invalid username
func TestInvalidUsername(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServer(t, &apiServer)
	defer test.CloseServer(ts)

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
	ts := initAPITestServerInvalid(t, &apiServer, "")
	defer test.CloseServer(ts)

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
	ts := initAPITestServer(t, &apiServer)
	defer test.CloseServer(ts)

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
	ts := initAPITestServer(t, &apiServer)
	defer test.CloseServer(ts)

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
	ts := initAPITestServerInvalid(t, &apiServer, "")
	defer test.CloseServer(ts)

	problems, err := apiServer.GetProblemListCPbook(3)
	if err == nil {
		t.Fatalf("Error %v", "expected end of JSON input")
	}
	if len(problems) != 0 {
		t.Fatalf("Error expected an empty response")
	}

	// empty JSON object
	ts = initAPITestServerInvalid(t, &apiServer, "{}")
	problems, err = apiServer.GetProblemListCPbook(3)
	if err == nil {
		t.Fatalf("Error %v", "expected json: cannot unmarshal object")
	}
	if len(problems) != 0 {
		t.Fatalf("Error expected an empty response")
	}
}

// Test get problems from CP book, empty number of problems
func TestGetCPBookProblemsEmpty(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServerInvalid(t, &apiServer, "[]")
	defer test.CloseServer(ts)

	problems, err := apiServer.GetProblemListCPbook(3)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if len(problems) != 0 {
		t.Fatalf("Error expected an empty response")
	}
}

// Test get problem list
func TestGetProblemList(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServer(t, &apiServer)
	defer test.CloseServer(ts)

	problems, err := apiServer.GetProblemList()
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if problems[100].Title != "The 3n + 1 problem" {
		t.Fatalf("Error: title does not match: %v", problems[100].Title)
	}
}

// Test get user submissions
func TestGetUserSubmissions(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServer(t, &apiServer)
	defer test.CloseServer(ts)

	submissions, err := apiServer.GetUserSubmissions(userid)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if submissions.Username != username ||
		len(submissions.Submissions) != nUserSubmissions {
		t.Fatalf("Error submissions do not match: read %d expected %d",
			len(submissions.Submissions), nUserSubmissions)
	}
}

// Test get problems from CP book, invalid response
func TestGetUserSubmissionsInvalidResponse(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServerInvalid(t, &apiServer, "")
	defer test.CloseServer(ts)

	submissions, err := apiServer.GetUserSubmissions(userid)
	if err == nil {
		t.Fatalf("Error %v", "expected json: cannot unmarshal object")
	}
	if submissions.Username == username ||
		len(submissions.Submissions) != 0 {
		t.Fatalf("Error expected an empty response")
	}

	// empty JSON object
	ts = initAPITestServerInvalid(t, &apiServer, "[]")

	submissions, err = apiServer.GetUserSubmissions(userid)
	if err == nil {
		t.Fatalf("Error %v", "expected json: cannot unmarshal object")
	}
	if submissions.Username == username ||
		len(submissions.Submissions) != 0 {
		t.Fatalf("Error expected an empty response")
	}
}

// Test get problems from CP book, empty number of problems
func TestGetUserSubmissionsEmpty(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServerInvalid(t, &apiServer, "{}")
	defer test.CloseServer(ts)

	submissions, err := apiServer.GetUserSubmissions(userid)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if submissions.Username == username ||
		len(submissions.Submissions) != 0 {
		t.Fatalf("Error expected an empty response")
	}
}

// Test get problem details by number
func TestGetProblemByNum(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServer(t, &apiServer)
	defer test.CloseServer(ts)

	problem, err := apiServer.GetProblemByNum(problemNumber)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if problem.ProblemID != problemID {
		t.Fatalf("Error %v", "problem does not match")
	}
}

// Test get problem details by number, invalid response
func TestGetProblemByNumInvalidResponse(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServerInvalid(t, &apiServer, "")
	defer test.CloseServer(ts)

	problem, err := apiServer.GetProblemByNum(problemNumber)
	if err == nil {
		t.Fatalf("Error %v", "expected json: cannot unmarshal object")
	}
	if problem.ProblemID == problemID {
		t.Fatalf("Error expected an empty response")
	}

	// empty JSON object
	ts = initAPITestServerInvalid(t, &apiServer, "[]")

	problem, err = apiServer.GetProblemByNum(problemNumber)
	if err == nil {
		t.Fatalf("Error %v", "expected json: cannot unmarshal object")
	}
	if problem.ProblemID == problemID {
		t.Fatalf("Error expected an empty response")
	}
}

// Test get problem details by number, empty number of problems
func TestGetProblemByNumEmpty(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServerInvalid(t, &apiServer, "{}")
	defer test.CloseServer(ts)

	problem, err := apiServer.GetProblemByNum(problemNumber)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if problem.ProblemID == problemID {
		t.Fatalf("Error expected an empty response")
	}
}

// Test get problem details by id
func TestGetProblemByID(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServer(t, &apiServer)
	defer test.CloseServer(ts)

	problem, err := apiServer.GetProblemByID(problemID)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if problem.ProblemNumber != problemNumber {
		t.Fatalf("Error %v", "problem does not match")
	}
}

// Test get problem details by ID, invalid response
func TestGetProblemByIDInvalidResponse(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServerInvalid(t, &apiServer, "")
	defer test.CloseServer(ts)

	problem, err := apiServer.GetProblemByID(problemID)
	if err == nil {
		t.Fatalf("Error %v", "expected json: cannot unmarshal object")
	}
	if problem.ProblemNumber == problemNumber {
		t.Fatalf("Error expected an empty response")
	}

	// empty JSON object
	ts = initAPITestServerInvalid(t, &apiServer, "[]")

	problem, err = apiServer.GetProblemByID(problemID)
	if err == nil {
		t.Fatalf("Error %v", "expected json: cannot unmarshal object")
	}
	if problem.ProblemNumber == problemNumber {
		t.Fatalf("Error expected an empty response")
	}
}

// Test get problem details by ID, empty number of problems
func TestGetProblemByIDEmpty(t *testing.T) {
	var apiServer APIServer
	ts := initAPITestServerInvalid(t, &apiServer, "{}")
	defer test.CloseServer(ts)

	problem, err := apiServer.GetProblemByID(problemID)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if problem.ProblemNumber == problemNumber {
		t.Fatalf("Error expected an empty response")
	}
}

// Initialize the test environment
func TestMain(m *testing.M) {
	flag.BoolVar(&realTest, "real", false, "Test with real uHunt API server")
	flag.Parse()
	os.Exit(m.Run())
}
