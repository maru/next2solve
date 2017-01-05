// Next problem to solve
// https://github.com/maru/next2solve
//
// Tests for problems.go functionality
//
package problems

import (
	"flag"
	"net/http/httptest"
	test "next2solve/testing"
	"os"
	"testing"
	"time"
)

//
const (
	nCP3BookProblems  = 1655
	nCP3BookTitles    = 213
	nUnsolvedProblems = 1337
	problemID         = 1998
	problemNumber     = 11057
	problemTitle      = "Exact Sum"
	chapterTitle      = "Problem Solving Paradigms"
	userID            = "46232"
	username          = "chicapi"
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

// Test initialize API server and create the caches.
func TestInitAPIServer(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	if ts.URL != apiServer.GetUrl() {
		t.Fatalf("Expected API server URL %s, got %s", ts.URL, apiServer.GetUrl())
	}

	if ts.URL != apiServer.GetUrl() {
		t.Fatalf("Expected API server URL %s, got %s", ts.URL, apiServer.GetUrl())
	}

	cacheNames := []string{"userid", "submissions", "problem"}
	for _, s := range cacheNames {
		if _, ok := cache[s]; !ok {
			t.Fatalf("Error cache: %v not found", s)
		}
	}
	if _, ok := cache["notvalid"]; ok {
		t.Fatalf("Error cache: %v not expected", "notvalid")
	}
}

// Test get chapter name of a problem. Problem list must be loaded in previous
// test
func TestGetChapter(t *testing.T) {
	p := ProblemInfo{problemID, problemNumber, problemTitle, 0, 0, 0, false}
	if p.GetChapter() != chapterTitle {
		t.Fatalf("Expected %s, got %s", chapterTitle, p.GetChapter())
	}
}

// Test load problem list from CP3 book, but API server sends empty or invalid
// response.
func TestLoadProblemListCP3EmptyInvalidResponse(t *testing.T) {
	ts := initAPITestServerInvalid(t, []string{"[]", ""})
	defer test.CloseServer(ts)

	if len(cpProblems) != 0 {
		t.Fatalf("Expected %d problems", 0)
	}
	if len(cpTitles) != 0 {
		t.Fatalf("Expected %d titles", 0)
	}
	if len(problemList) != 0 {
		t.Fatalf("Expected %d problems, got %d", 0, len(problemList))
	}

	// Invalid response
	loadProblemListCP3()

	if len(cpProblems) != 0 {
		t.Fatalf("Expected %d problems", 0)
	}
	if len(cpTitles) != 0 {
		t.Fatalf("Expected %d titles", 0)
	}
	if len(problemList) != 0 {
		t.Fatalf("Expected %d problems, got %d", 0, len(problemList))
	}
}

// Test the number of titles and problems loaded from CP3 book.
func TestLoadProblemListCP3(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	if len(cpProblems) != nCP3BookProblems {
		t.Fatalf("Expected %d problems", nCP3BookProblems)
	}
	if len(cpTitles) != nCP3BookTitles {
		t.Fatalf("Expected %d titles", nCP3BookTitles)
	}
	if len(problemList) != nCP3BookProblems {
		t.Fatalf("Expected %d problems, got %d", nCP3BookProblems, len(problemList))
	}
}

func TestRefreshProblemCache(t *testing.T) {
	ts := test.InitAPITestServer(t)
	defer test.CloseServer(ts)

	// InitAPIServer
	apiServer.Init(ts.URL)
	// Create cache for each type of object
	cache = make(map[string]*Cache)
	cache["userid"] = NewCache(cacheDurationUser)
	cache["submissions"] = NewCache(cacheDurationSubmissions)
	cache["problem"] = NewCache(time.Second)

	// Load list of problems to solve from the CP3 book
	loadProblemListCP3()

	go refreshProblemCache(time.Second)
	time.Sleep(3 * time.Second)
	// Quit refreshProblemCache
	quitRefreshCache <- true
}

// Test get problem information, invalid problem ID
func TestGetProblemInvalid(t *testing.T) {
	// API server sends empty JSON object
	ts := initAPITestServerInvalid(t, []string{"[]", "{}", "[]", ""})
	defer test.CloseServer(ts)

	problem := getProblem(problemID + 1000000)
	if problem.Number == problemNumber || problem.Title == problemTitle {
		t.Fatalf("Error expected empty problem information")
	}

	// Invalid JSON
	problem = getProblem(problemID + 1000000)
	if problem.Number == problemNumber || problem.Title == problemTitle {
		t.Fatalf("Error expected empty problem information")
	}
}

// Test get problem information, valid problem ID
func TestGetProblemOk(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problem := getProblem(problemID)
	if problem.Number != problemNumber || problem.Title != problemTitle {
		t.Fatalf("Error problem number or title does not match")
	}
}

// Test get problem information, valid problem ID, not found in chache.
func TestGetProblemSetCache(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problemID := 70
	problemNumber := 134
	problemTitle := "Loglan-A Logical Language"
	problem := getProblem(problemID)
	if problem.Number != problemNumber || problem.Title != problemTitle {
		t.Fatalf("Error problem number or title does not match")
	}
}

// Test get userID with invalid username
func TestGetUserIDInvalid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	invalidUsername := "not_" + username
	id, err := GetUserID(invalidUsername)
	if err == nil || err.Error() != "Username not found" {
		t.Fatalf("Expected error")
	}
	if id != "" {
		t.Fatalf("Expected empty userID")
	}
}

// Test get userID with valid username
func TestGetUserIDValid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	id, err := GetUserID(username)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	if id != userID {
		t.Fatalf("Expected userID %v got %v", userID, id)
	}
}

// Test getUserSubmissions, submissions cached
func TestGetUserSubmissions(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	sub1 := getUserSubmissions(userID)
	sub2 := getUserSubmissions(userID)
	if len(sub1) != len(sub2) && sub1[0].SubmissionID != sub2[0].SubmissionID {
		t.Fatalf("Error submissions not equal")
	}
}

// Test GetUnsolvedProblemsCPBook
func TestGetUnsolvedProblemsCPBookOk(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblemsCPBook(userID, "star")
	if len(problems) != nUnsolvedProblems {
		t.Fatalf("Expected %d problems to solve, got %d", nUnsolvedProblems, len(problems))
	}
}

// Test GetUnsolvedProblemsCPBook, invalid API response
func TestGetUnsolvedProblemsCPBookInvalidResponse(t *testing.T) {
	ts := initAPITestServerInvalid(t, []string{"[]", ""})
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblemsCPBook(userID, "star")
	if len(problems) != 0 {
		t.Fatalf("Expected empty problem list")
	}
}

// Test GetUnsolvedProblems, valid userID
func TestGetUnsolvedProblemsValidUserid(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblems(userID, "star")
	if len(problems) != nUnsolvedProblems {
		t.Fatalf("Expected %d problems to solve, got %d", nUnsolvedProblems, len(problems))
	}
}

// Test GetUnsolvedProblemRandom
func TestGetUnsolvedProblemRandomOk(t *testing.T) {
	ts := initAPITestServer(t)
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblemRandom(userID)
	if len(problems) != 1 {
		t.Fatalf("Expected problem list")
	}
}

// Test GetUnsolvedProblemRandom, invalid API response
func TestGetUnsolvedProblemRandomInvalidResponse(t *testing.T) {
	ts := initAPITestServerInvalid(t, []string{"[]", ""})
	defer test.CloseServer(ts)

	problems := GetUnsolvedProblemRandom(userID)
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
