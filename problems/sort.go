// Next problem to solve
// https://github.com/maru/next2solve
//
// Sort problems
//

package problems

import (
	"sort"
)

type ProblemList []int

// Return problem list length
func (a ProblemList) Len() int {
	return len(a)
}

// Swap two elements of problem list
func (a ProblemList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ProblemList) Less(i, j int) bool {
	pID := a[i]
	qID := a[j]
	return pID < qID
}

type ProblemListStar ProblemList

// Return problem list length
func (a ProblemListStar) Len() int {
	return len(a)
}

// Swap two elements of problem list
func (a ProblemListStar) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Compare two elements of problem list:
// Sort ProblemListStar by star first, level asc, acratio desc, dacu desc
func (a ProblemListStar) Less(i, j int) bool {
	pID := a[i]
	qID := a[j]
	p := getProblem(pID)
	q := getProblem(qID)
	cpP := cpProblems[pID]
	cpQ := cpProblems[qID]

	if p.Star != q.Star {
		return p.Star && !q.Star
	}
	if p.Level != q.Level {
		return p.Level < q.Level
	}
	if cpP.Chapter != cpQ.Chapter {
		return cpP.Chapter < cpQ.Chapter
	}
	if p.AcRatio != q.AcRatio {
		return p.AcRatio > q.AcRatio
	}
	if p.Dacu != q.Dacu {
		return p.Dacu > q.Dacu
	}
	return pID < qID
}

type ProblemListCategory ProblemList

// Return problem list length
func (a ProblemListCategory) Len() int {
	return len(a)
}

// Swap two elements of problem list
func (a ProblemListCategory) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Compare two elements of problem list:
// Sort ProblemListCategory by chapter, level asc, star, acratio desc, dacu desc
func (a ProblemListCategory) Less(i, j int) bool {
	pID := a[i]
	qID := a[j]
	p := getProblem(pID)
	q := getProblem(qID)
	cpP := cpProblems[pID]
	cpQ := cpProblems[qID]

	if cpP.Chapter != cpQ.Chapter {
		return cpP.Chapter < cpQ.Chapter
	}
	if p.Level != q.Level {
		return p.Level < q.Level
	}
	if p.Star != q.Star {
		return p.Star && !q.Star
	}
	if p.Dacu != q.Dacu {
		return p.Dacu > q.Dacu
	}
	return pID < qID
}

type ProblemListLevel ProblemList

// Return problem list length
func (a ProblemListLevel) Len() int {
	return len(a)
}

// Swap two elements of problem list
func (a ProblemListLevel) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Compare two elements of problem list:
// Sort ProblemListCategory by level asc, star, chapter, acratio desc, dacu desc
func (a ProblemListLevel) Less(i, j int) bool {
	pID := a[i]
	qID := a[j]
	p := getProblem(pID)
	q := getProblem(qID)

	if p.Level != q.Level {
		return p.Level < q.Level
	}
	if p.Dacu != q.Dacu {
		return p.Dacu > q.Dacu
	}
	return pID < qID
}

// Sort problem list
func sortProblemList(problemList []int, orderBy string) {
	switch orderBy {
	case "star":
		sort.Sort(ProblemListStar(problemList))
	case "cat":
		sort.Sort(ProblemListCategory(problemList))
	case "lev":
		sort.Sort(ProblemListLevel(problemList))
	}
}
