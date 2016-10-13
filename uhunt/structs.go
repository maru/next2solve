// Next problem to solve
// https://github.com/maru/next2solve
//
// uHunt API structs

package uhunt

import (
	"math"
)

const (
	VerdictSubmissionError    = 10
	VerdictCantBeJudged       = 15
	VerdictInQueue            = 20
	VerdictCompileError       = 30
	VerdictRestrictedFunction = 35
	VerdictRuntimeError       = 40
	VerdictOutputLimit        = 45
	VerdictTimeLimit          = 50
	VerdictMemoryLimit        = 60
	VerdictWrongAnswer        = 70
	VerdictPresentationError  = 80
	VerdictAccepted           = 90
)

type APIServer struct {
	urlServer string
}

type APICPBookChapter struct {
	Title       string                `json:"title"`
	Subchapters []APICPBookSubchapter `json:"arr"`
}
type APICPBookSubchapter struct {
	Title    string        `json:"title"`
	Sections []interface{} `json:"arr"`
}

type APIProblem struct {
	ProblemID              int    `json:"pid"`
	ProblemNumber          int    `json:"num"`
	Title                  string `json:"title"`
	Dacu                   int    `json:"dacu"`
	BestRuntime            int    `json:"mrun"`
	BestUsedMemory         int    `json:"mmem"`
	NumNoVerdict           int    `json:"nover"`
	NumSubmissionError     int    `json:"sube"`
	NumCantBeJudged        int    `json:"noj"`
	NumInQueue             int    `json:"inq"`
	NumCompilationError    int    `json:"ce"`
	NumRestrictedFunction  int    `json:"rf"`
	NumRuntimeError        int    `json:"re"`
	NumOutputLimitExceeded int    `json:"ole"`
	NumTimeLimitExceeded   int    `json:"tle"`
	NumMemoryLimitExceeded int    `json:"mle"`
	NumWrongAnswer         int    `json:"wa"`
	NumPresentationError   int    `json:"pe"`
	NumAccepted            int    `json:"ac"`
	RunTimeLimit           int    `json:"rtl"`
	Status                 int    `json:"status"`
	Rej                    int    `json:"rej"`
}

// Get level (value between 1 and 10)
func (p *APIProblem) GetLevel() int {
	dacuLog := math.Log(float64(p.Dacu) + 1)
	maxLevel := math.Floor(math.Min(10, dacuLog))
	level := math.Max(1, 10-maxLevel)
	ret := int(level)
	return ret
}

// Get acceptance ratio
func (p *APIProblem) GetAcceptanceRatio() int {
	total := p.GetTotalSubmissions()
	if total <= 0 {
		return 0
	}
	return (p.NumAccepted * 100.0) / total
}

// Get total number of submissions
func (p *APIProblem) GetTotalSubmissions() int {
	return p.NumNoVerdict + p.NumSubmissionError + p.NumCantBeJudged + p.NumInQueue +
		p.NumCompilationError + p.NumRestrictedFunction + p.NumRuntimeError +
		p.NumOutputLimitExceeded + p.NumTimeLimitExceeded + p.NumMemoryLimitExceeded +
		p.NumWrongAnswer + p.NumPresentationError + p.NumAccepted
}

type APIUserSubmissions struct {
	Name        string          `json:"name"`
	Username    string          `json:"uname"`
	TmpSubs     [][]int         `json:"subs"`
	Submissions []APISubmission `json:"-"`
}
type APISubmission struct {
	SubmissionID   int
	ProblemID      int
	VerdictID      int
	Runtime        int
	Time           int
	Language       int
	SubmissionRank int
}

func (s *APISubmission) IsAccepted() bool {
	return s.VerdictID == VerdictAccepted
}
