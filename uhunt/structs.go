// Next problem to solve
// https://github.com/maru/next2solve
//
// uHunt API structs

package uhunt

import (
	"math"
)

type APIServer struct {
	urlServer string
}

type APICPBookChapter struct {
	Title       string             `json:"title"`
	Subchapters []APICPBookSubchapter `json:"arr"`
}
type APICPBookSubchapter struct {
	Title          string        `json:"title"`
	Sections []interface{} `json:"arr"`
}

type APIProblem struct {
	ProblemID              int64  `json:"pid"`
	ProblemNumber          int64  `json:"num"`
	Title                  string `json:"title"`
	Dacu                   int64  `json:"dacu"`
	BestRuntime            int64  `json:"mrun"`
	BestUsedMemory         int64  `json:"mmem"`
	NumNoVerdict           int64  `json:"nover"`
	NumSubmissionError     int64  `json:"sube"`
	NumCantBeJudged        int64  `json:"noj"`
	NumInQueue             int64  `json:"inq"`
	NumCompilationError    int64  `json:"ce"`
	NumRestrictedFunction  int64  `json:"rf"`
	NumRuntimeError        int64  `json:"re"`
	NumOutputLimitExceeded int64  `json:"ole"`
	NumTimeLimitExceeded   int64  `json:"tle"`
	NumMemoryLimitExceeded int64  `json:"mle"`
	NumWrongAnswer         int64  `json:"wa"`
	NumPresentationError   int64  `json:"pe"`
	NumAccepted            int64  `json:"ac"`
	RunTimeLimit           int64  `json:"rtl"`
	Status                 int64  `json:"status"`
	Rej                    int64  `json:"rej"`
}

// Get level (value between 1 and 10)
func (p *APIProblem) GetLevel() int64 {
	return int64(math.Max(1, 10-math.Floor(math.Min(10, math.Log(float64(p.Dacu))))))
}

// Get acceptance ratio
func (p *APIProblem) GetAcceptanceRatio() int64 {
	return p.NumAccepted * 100.0 / p.GetTotalSubmissions()
}

// Get total number of submissions
func (p *APIProblem) GetTotalSubmissions() int64 {
	return p.NumNoVerdict + p.NumSubmissionError + p.NumCantBeJudged + p.NumInQueue +
		p.NumCompilationError + p.NumRestrictedFunction + p.NumRuntimeError +
		p.NumOutputLimitExceeded + p.NumTimeLimitExceeded + p.NumMemoryLimitExceeded +
		p.NumWrongAnswer + p.NumPresentationError + p.NumAccepted
}

type APIUserSubmissions struct {
	Name        string      `json:"name"`
	Username    string      `json:"uname"`
	Submissions [][]int64 `json:"subs"`
}
