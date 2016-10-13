// Next problem to solve
// https://github.com/maru/next2solve
//
// HTTP handler

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"next2solve/problems"
)

type TemplateData struct {
	UsernameError string
	UserID        string
	Username      string
	Problems      []problems.ProblemInfo
}

var (
	funcMap = template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}
	templates = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/header.html",
		"templates/footer.html", "templates/index.html", "templates/lucky.html",
		"templates/problems.html"))
)

// Render page using a template with data
func renderPage(w http.ResponseWriter, tmpl string, data interface{}) {
	templates = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/header.html",
		"templates/footer.html", "templates/index.html", "templates/lucky.html",
		"templates/problems.html"))
	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		fmt.Fprintf(w, "Error %v", err)
	}
}

// Show unsolved problems
func showProblems(w http.ResponseWriter, data TemplateData) {
	data.Problems = problems.GetUnsolvedProblems(data.UserID)
	if len(data.Problems) == 0 {
		data = TemplateData{"No problems to solve", "", data.Username, nil}
		renderPage(w, "index", data)
		return
	}
	renderPage(w, "problems", data)
}

// Show a random unsolved problem
func showRandomProblem(w http.ResponseWriter, data TemplateData) {
	// Choose a problem with lowest dacu, starred first
	data.Problems = problems.GetUnsolvedProblemRandom(data.UserID)
	if len(data.Problems) == 0 {
		data = TemplateData{"No problems to solve", "", data.Username, nil}
		renderPage(w, "index", data)
		return
	}
	renderPage(w, "lucky", data)
}

// Get user information from cookies
func getTemplateData(r *http.Request) TemplateData {
	data := TemplateData{}
	if cookie, err := r.Cookie("userid"); err == nil {
		data.UserID = cookie.Value
	}
	if cookie, err := r.Cookie("username"); err == nil {
		data.Username = cookie.Value
	}
	return data
}

// Set user information in cookies
func setTemplateData(w http.ResponseWriter, data *TemplateData, userid, username string) {
	if data.UserID != userid {
		cookie := http.Cookie{Name: "userid", Value: userid}
		http.SetCookie(w, &cookie)
		data.UserID = userid
	}
	if data.Username != username {
		cookie := http.Cookie{Name: "username", Value: username}
		http.SetCookie(w, &cookie)
		data.Username = username
	}
}

// Handles requests
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	data := getTemplateData(r)
	if r.Method == "POST" {
		// Show problems to solve
		username := r.PostFormValue("username")
		// Check if username is valid
		userid, err := problems.GetUserID(username)
		if err != nil {
			data = TemplateData{err.Error(), "", username, nil}
			renderPage(w, "index", data)
			return
		}
		// Set user information in a cookie
		setTemplateData(w, &data, userid, username)

		// Show all unsolved problems
		if r.PostFormValue("show-problems") != "" {
			// Show all unsolved problems
			showProblems(w, data)
			return
		} else {
			// Show a random unsolved problem
			showRandomProblem(w, data)
			return
		}
	}

	// GET - Default
	renderPage(w, "index", data)
}

// Set handlers, initialize API server and start HTTP server
func HttpServerStart(addr string, apiUrl string) {
	problems.InitAPIServer(apiUrl)
	http.HandleFunc("/", RequestHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
