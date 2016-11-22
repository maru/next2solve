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

// Compare two elements of problem list:
// Sort problemList by star first, level asc, acratio desc, dacu desc
func (a ProblemList) Less(i, j int) bool {
	pID := a[i]
	qID := a[j]
	p := getProblem(pID)
	q := getProblem(qID)
	cpP := cpProblems[pID]
	cpQ := cpProblems[qID]
	if p.Star && !q.Star {
		return true
	}
	if !p.Star && q.Star {
		return false
	}
	if p.Level < q.Level {
		return true
	}
	if p.Level > q.Level {
		return false
	}
	if cpP.Chapter < cpQ.Chapter {
		return true
	}
	if cpP.Chapter > cpQ.Chapter {
		return false
	}
	if p.AcRatio > q.AcRatio {
		return true
	}
	if p.AcRatio < q.AcRatio {
		return false
	}
	if p.Dacu > q.Dacu {
		return true
	}
	if p.Dacu < q.Dacu {
		return false
	}
	return a[i] < a[j]
}

// Sort problem list
func sortProblemList() {
	sort.Sort(ProblemList(problemList))
}
