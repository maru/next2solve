// Next problem to solve
// https://github.com/maru/next2solve
//
// Problems
//
// Number of Distinct Accepted User (DACU)

package main

import (
	"errors"
	"math"
	"math/rand"
	"next2solve/uhunt"
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
}

var (
	apiServer  uhunt.APIServer
	cpProblems map[int64]CPProblem
)

func InitAPIServer(url string) {
	apiServer.Init(url)
}

// Call the API to get the user id from the username
func GetUserID(username string) (string, error) {
	id, _ := apiServer.GetUserID(username)
	if id == "0" {
		return "", errors.New("Username not found")
	}
	return id, nil
}

// Call the API to get the problem list (to solve) and the solved problems by
// the user.
// Get the unsolved problems, sort by level and acceptance ratio (desc).
func GetUnsolvedProblems(userid string) []ProblemInfo {
	cpBook, err := apiServer.GetProblemList()
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

	userSubs, err := apiServer.GetUserSubmissions(userid)
	if err != nil {
		return []ProblemInfo{}
	}

	// Get only accepted (distinct) problems
	userProblems := make(map[int64]bool)
	for _, p := range userSubs.Submissions {
		if p.VerdictID == uhunt.VerdictAccepted {
			userProblems[int64(p.ProblemID)] = true
		}
	}

	// Filter solved problems
	var unsolved []ProblemInfo
	for _, pnum := range problems {
		p, _ := apiServer.GetProblemByNum(pnum)
		if _, ok := userProblems[pnum]; !ok {
			unsolved = append(unsolved, ProblemInfo{pnum, p.ProblemNumber, p.Title,
				p.GetLevel(), p.GetAcceptanceRatio()})
		}
	}
	return unsolved
}

// Get the unsolved problems and return one random problem.
func GetUnsolvedProblemRandom(userid string) []ProblemInfo {
	// Choose a problem with lowest dacu, starred first
	unsolved := GetUnsolvedProblems(userid)
	if len(unsolved) > 0 {
		r := rand.Intn(len(unsolved))
		return []ProblemInfo{unsolved[r]}
	}
	return []ProblemInfo{}
}
