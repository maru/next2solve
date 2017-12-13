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
// Sort ProblemListLevel by level asc, dacu desc, id asc
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

type ProblemListSubmissions ProblemList

// Return problem list length
func (a ProblemListSubmissions) Len() int {
	return len(a)
}

// Swap two elements of problem list
func (a ProblemListSubmissions) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Compare two elements of problem list:
// Sort ProblemListSubmissions by submissions desc, id asc
func (a ProblemListSubmissions) Less(i, j int) bool {
	pID := a[i]
	qID := a[j]
	p := getProblem(pID)
	q := getProblem(qID)

	if p.TotalSubmissions != q.TotalSubmissions {
		return p.TotalSubmissions > q.TotalSubmissions
	}
	return pID < qID
}

type ProblemListAccepted ProblemList

// Return problem list length
func (a ProblemListAccepted) Len() int {
	return len(a)
}

// Swap two elements of problem list
func (a ProblemListAccepted) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Compare two elements of problem list:
// Sort ProblemListAccepted by accepted desc, id asc
func (a ProblemListAccepted) Less(i, j int) bool {
	pID := a[i]
	qID := a[j]
	p := getProblem(pID)
	q := getProblem(qID)

	if p.TotalAccepted != q.TotalAccepted {
		return p.TotalAccepted > q.TotalAccepted
	}
	return pID < qID
}

type ProblemListACRatio ProblemList

// Return problem list length
func (a ProblemListACRatio) Len() int {
	return len(a)
}

// Swap two elements of problem list
func (a ProblemListACRatio) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Compare two elements of problem list:
// Sort ProblemListACRatio by AcRatio desc, id asc
func (a ProblemListACRatio) Less(i, j int) bool {
	pID := a[i]
	qID := a[j]
	p := getProblem(pID)
	q := getProblem(qID)

	if p.AcRatio != q.AcRatio {
		return p.AcRatio > q.AcRatio
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
  case "sub":
		sort.Sort(ProblemListSubmissions(problemList))
  case "ac":
		sort.Sort(ProblemListAccepted(problemList))
  case "acr":
		sort.Sort(ProblemListACRatio(problemList))
  case "dacu":
		sort.Sort(ProblemListLevel(problemList))
	}
}
