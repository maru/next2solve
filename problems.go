// Next problem to solve
// https://github.com/maru/next2solve
//
// Problems
//
// Number of Distinct Accepted User (DACU)

package main

import (
	// "next2solve/uhunt"
	// "math"
	"math/rand"
)

type CPProblem struct {
	Star bool
	Chapter int
	Subchapter int
	Section int
}

type ProblemInfo struct {
	ProblemID      int
	ProblemNumber  int64
	ProblemTitle   string
	ProblemLevel   int64
	ProblemAcRatio int64
}

var (
	apiServer uhunt.APIServer
	cpProblems map[int]CPProblem
)

func InitAPIServer(url string) {
	apiServer.Init(url)
}

// Call the API to get the user id from the username
func GetUserID(username string) (string, error) {
	id, err := apiServer.GetUserID(invalidUsername)
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
	var problems []int
	cpProblems = make(map[int]CPProblem)
	numChapters := 0
	numSubChapters := 100
	numSubSubChapters := 1000
	for _, chapter := range cpBook {
		numChapters++
		for _, subchapter := range chapter.Subchapters {
			numSubChapters++
			for _, subsubchapter := range subchapter.Subsubchapters {
				numSubSubChapters++
				arr := subsubchapter.([]interface{})
				for _, p := range arr[1:] {
					pid := int(math.Abs(p.(float64)))
					problems = append(problems, pid)
					cpProblems[pid] = CPProblem{p.(float64) < 0, numChapters, numSubChapters, numSubSubChapters}
				}
			}
		}
	}

	userSubs, err := apiServer.GetUserSubmissions(userid)
	if err != nil {
		return []ProblemInfo{}
	}

	// Get only accepted (distinct) problems
	// An element in the array contains:
	//  0   Submission ID
	//  1   Problem ID  (* we want this)
	//  2   Verdict ID  (* and this with value 90 : Accepted)
	//  3   Runtime
	//  4   Submission Time (unix timestamp)
	//  5   Language ID (1=ANSI C, 2=Java, 3=C++, 4=Pascal, 5=C++11)
	//  6   Submission Rank
	userProblems := make(map[int]bool)
	for _, p := range userSubs.Submissions {
		pid := int(math.Abs(p[1]))
		if p[2] == 90 {
			userProblems[pid] = true
		}
	}

	println("problems", len(problems))
	println("userProblems", len(userProblems))
	// Filter solved problems
	var unsolved []ProblemInfo
	for _, pnum := range problems {
		p, _ := apiServer.GetProblemInfoByNum(pnum)
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
