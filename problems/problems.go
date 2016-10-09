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
	"os"
)

type CPProblem struct {
	Star       bool
	Chapter    int
	Subchapter int
	Section    int
}

type ProblemInfo struct {
	ID      int64
	Number  int64
	Title   string
	Level   int64
	AcRatio int64
	Dacu    int64
}

var (
	apiServer  uhunt.APIServer
	cpProblems map[int64]CPProblem
	cache *Cache
)

func InitAPIServer(url string) {
	apiServer.Init(url)
	cache = NewCache()
	if err := cache.CreateNamespace("userid"); err != nil {
		println(err.Error())
		os.Exit(1)
	}
	if err := cache.CreateNamespace("submissions"); err != nil {
		println(err.Error())
		os.Exit(1)
	}
	if err := cache.CreateNamespace("problem"); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
	// cache.LoadProblemListCP3()

// Call the API to get the user id from the username
func GetUserID(username string) (string, error) {
	// Get userid from cache, if found and valid.
	// Otherwise, call the API and set the value in the cache
	id, ok := cache.Get("userid", username)
	if !ok {
		id, _ = apiServer.GetUserID(username)
		cache.Set("userid", username, id)
	}
	// Check userid
	if id.(string) == "0" {
		return "", errors.New("Username not found")
	}
	return id.(string), nil
}

// Get the unsolved problems, sort by level and acceptance ratio (desc).
// Calls the API to get the problem list (from the CP3 book), the details of
// each problem and the submissions by the user.
func GetUnsolvedProblemsCPBook(userid string) []ProblemInfo {
	// var userSubs APIUserSubmissions
	// if cache.GetUserSubmissions(userid) {
	//
	// }
	userSubs, err := apiServer.GetUserSubmissions(userid)
	if err != nil || userSubs.Username == "" {
		return []ProblemInfo{}
	}

	// Get only accepted (distinct) problems
	userProblems := make(map[int64]bool)
	for _, p := range userSubs.Submissions {
		if p.VerdictID == uhunt.VerdictAccepted {
			userProblems[int64(p.ProblemID)] = true
		}
	}

	// Get problem list of CP3 book
	cpBook, err := apiServer.GetProblemListCPbook(3)
	if err != nil {
		return []ProblemInfo{}
	}
	// Create an array with the problem ids
	var problems []int64
	cpProblems = make(map[int64]CPProblem)
	numChapters := 0
	numSubchapters := 100
	numSections := 1000
	for _, chapter := range cpBook {
		numChapters++
		for _, subchapter := range chapter.Subchapters {
			numSubchapters++
			for _, subsubchapter := range subchapter.Sections {
				numSections++
				arr := subsubchapter.([]interface{})
				for _, p := range arr[1:] {
					pid := int64(math.Abs(p.(float64)))
					problems = append(problems, pid)
					cpProblems[pid] = CPProblem{p.(float64) < 0, numChapters, numSubchapters, numSections}
				}
			}
		}
	}

	// Filter solved problems
	var unsolved []ProblemInfo
	for _, pnum := range problems {
		p, _ := apiServer.GetProblemByNum(pnum)
		if _, ok := userProblems[pnum]; !ok {
			unsolved = append(unsolved, ProblemInfo{pnum, p.ProblemNumber, p.Title,
				p.GetLevel(), p.GetAcceptanceRatio(), p.Dacu})
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
