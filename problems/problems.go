// Next problem to solve
// https://github.com/maru/next2solve
//
// Problems
//

package problems

import (
	"errors"
	"math"
	"math/rand"
	"next2solve/uhunt"
	"time"
)

type CPProblem struct {
	Star       bool
	Chapter    int
	Subchapter int
	Section    int
}

type ProblemInfo struct {
	ID      int
	Number  int
	Title   string
	Level   int
	AcRatio int
	Dacu    int
	// CPInfo  CPProblem
}

const (
	cacheDurationUser        = time.Hour
	cacheDurationSubmissions = time.Hour
	cacheDurationProblem     = 15 * time.Minute
)

var (
	apiServer   uhunt.APIServer
	cache       map[string]*Cache
	cpProblems  map[int]CPProblem
	cpTitles    map[int]string
	problemList []int
)

// Initialize API server and cache
func InitAPIServer(url string) {
	// Set API sever URL
	apiServer.Init(url)

	// Create cache for each type of object
	cache = make(map[string]*Cache)
	cache["userid"] = NewCache(cacheDurationUser)
	cache["submissions"] = NewCache(cacheDurationSubmissions)
	cache["problem"] = NewCache(cacheDurationProblem)

	// Load list of problems to solve from the CP3 book
	loadProblemListCP3()
	// Start Problem cache refresh in background
	go refreshProblemCache(cacheDurationProblem-time.Minute)
}

// Load chapter titles and the list of problems to solve from the CP3 book.
func loadProblemListCP3() {
	println("Loading problems...")
	// Get problem list of CP3 book
	cpBook, err := apiServer.GetProblemListCPbook(3)
	if err != nil {
		println("Error: couldn't load CP3 problem list from API")
		return
	}
	// Initialize
	cpProblems = make(map[int]CPProblem)
	cpTitles = make(map[int]string)
	problemList = []int{}
	numChapter := 0
	numSubchapter := 100
	numSection := 1000

	// Load titles and problems
	for _, chapter := range cpBook {
		// Chapter
		numChapter++
		cpTitles[numChapter] = chapter.Title
		for _, subchapter := range chapter.Subchapters {
			// Subchapter
			numSubchapter++
			cpTitles[numSubchapter] = subchapter.Title
			for _, section := range subchapter.Sections {
				// Section
				numSection++
				arr := section.([]interface{})
				cpTitles[numSection] = arr[0].(string)
				for _, problemNumber := range arr[1:] {
					// Get problem from API server
					pNum := int(math.Abs(problemNumber.(float64)))
					p, err := apiServer.GetProblemByNum(pNum)
					if err != nil {
						println("Error: couldn't load problem ", pNum, ":", err.Error())
						continue
					}
					// Set problem in cache
					problem := ProblemInfo{p.ProblemID, p.ProblemNumber, p.Title,
						p.GetLevel(), p.GetAcceptanceRatio(), p.Dacu}
					pID := p.ProblemID
					cache["problem"].Set(string(pID), problem)

					// Save CP3 problem information
					if _, ok := cpProblems[pID]; !ok {
						problemList = append(problemList, pID)
						cpProblems[pID] = CPProblem{problemNumber.(float64) < 0,
							numChapter, numSubchapter, numSection}
					}
				}
			}
		}
	}
	// Sort problemList by star first, level asc, acratio desc, dacu desc
	println("Done.")

}

// Refresh problem cache in background.
func refreshProblemCache(duration time.Duration) {
	for ; ; {
		timer1 := time.NewTimer(duration)
	  <-timer1.C
		println("refreshProblemCache")
		for _, pID := range problemList {
			getProblem(pID)
		}
	}
}

// Return problem information from cache first, otherwise from API server.
// If any error occurs, return empty problem.
func getProblem(pID int) ProblemInfo {
	problem, ok := cache["problem"].Get(string(pID))
	if !ok {
		p, err := apiServer.GetProblemByID(pID)
		if err != nil {
			println("Error: couldn't load problem ID", pID, ":", err.Error())
			return ProblemInfo{}
		}
		// Set problem in cache
		problem = ProblemInfo{p.ProblemID, p.ProblemNumber, p.Title,
			p.GetLevel(), p.GetAcceptanceRatio(), p.Dacu}
		pID := p.ProblemID
		cache["problem"].Set(string(pID), problem)
	}
	return problem.(ProblemInfo)
}

// Call the API to get the user id from the username.
func GetUserID(username string) (string, error) {
	// Get userid from cache, if found and valid.
	// Otherwise, call the API and set the value in the cache
	id, ok := cache["userid"].Get(username)
	if !ok {
		id, _ = apiServer.GetUserID(username)
		cache["userid"].Set(username, id)
	}
	// Check userid
	if id.(string) == "0" {
		return "", errors.New("Username not found")
	}
	return id.(string), nil
}

// Get user submissions from cache first, otherwise from API server.
// If any error occurs, return empty array.
func getUserSubmissions(userid string) []uhunt.APISubmission {
	if userSubs, ok := cache["submissions"].Get(userid); ok {
		us := userSubs.(uhunt.APIUserSubmissions)
		return us.Submissions
	}
	userSubs, err := apiServer.GetUserSubmissions(userid)
	if err != nil || userSubs.Username == "" {
		return []uhunt.APISubmission{}
	}
	cache["submissions"].Set(userid, userSubs)
	return userSubs.Submissions
}

// Get the unsolved problems, sort by level and acceptance ratio (desc).
// Calls the API to get the problem list (from the CP3 book), the details of
// each problem and the submissions by the user.
func GetUnsolvedProblemsCPBook(userid string) []ProblemInfo {

	// Get only accepted (distinct) problems
	userProblems := make(map[int]bool)
	submissions := getUserSubmissions(userid)
	for _, subm := range submissions {
		if subm.IsAccepted() {
			userProblems[subm.ProblemID] = true
		}
	}
	// Filter solved problems
	var unsolved []ProblemInfo
	for _, pID := range problemList {
		if _, ok := userProblems[pID]; !ok {
			unsolved = append(unsolved, getProblem(pID))
		}
	}
	return unsolved
}

func GetUnsolvedProblems(userid string) []ProblemInfo {
	return GetUnsolvedProblemsCPBook(userid)
}

// Get the unsolved problems by GetUnsolvedProblems and return one random problem.
func GetUnsolvedProblemRandom(userid string) []ProblemInfo {
	// Choose a problem with lowest dacu, starred first
	unsolved := GetUnsolvedProblems(userid)
	if len(unsolved) > 0 {
		r := rand.Intn(len(unsolved))
		return []ProblemInfo{unsolved[r]}
	}
	return []ProblemInfo{}
}
