// Next problem to solve
// https://github.com/maru/next2solve
//
// Problems
//
// Number of Distinct Accepted User (DACU)

package main

import (
	// "fmt"
	"math"
	// "net/http"
)

type ProblemInfo struct {
	ProblemID      string
	ProblemNumber  string
	ProblemTitle   string
	ProblemLevel   int64
	ProblemAcRatio float64
}

func getLevel(dacu int64) int64 {
	return int64(math.Max(1, 10-math.Floor(math.Min(10, math.Log(float64(dacu))))))
}

//
//
//
func getUnsolvedProblems(userid string) []ProblemInfo {
	var unsolved []ProblemInfo
	// problems := apiGetProblemList()
	// userProblems := apiGetUserProblems(userid)
	unsolved = append(unsolved, ProblemInfo{"15143", "31415", "Problem Title 1", 0, 0})
	unsolved = append(unsolved, ProblemInfo{"15143", "31415", "Problem Title 2", 1, 2})
	return unsolved
}

//
//
//
func getUnsolvedProblemRandom(userid string) []ProblemInfo {
	// Choose a problem with lowest dacu, starred first
	var unsolved []ProblemInfo
	// problems := apiGetProblemList()
	// userProblems := apiGetUserProblems(userid)
	unsolved = append(unsolved, ProblemInfo{"15143", "31415", "Problem Title", 0, 0})
	return unsolved
}
