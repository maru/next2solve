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

type ProblemInfo struct {
	ProblemID      int
	ProblemNumber  int64
	ProblemTitle   string
	ProblemLevel   int64
	ProblemAcRatio int64
}

// Get unsolved problems for a user, sort by level and acceptance ratio (desc)
func GetUnsolvedProblems(userid string) []ProblemInfo {
	problems, err := apiServer.GetProblemList()
	if err != nil {
		return []ProblemInfo{}
	}
	userProblems, err := apiServer.GetUserProblems(userid)
	if err != nil {
		return []ProblemInfo{}
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

//
//
//
func GetUnsolvedProblemRandom(userid string) []ProblemInfo {
	// Choose a problem with lowest dacu, starred first
	unsolved := GetUnsolvedProblems(userid)
	r := rand.Intn(len(unsolved))
	return []ProblemInfo{unsolved[r]}
}
