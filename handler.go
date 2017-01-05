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
	UsernameError   string
	UserID          string
	Username        string
	Problems        []problems.ProblemInfo
	IsOrderStar     bool
	IsOrderCategory bool
	IsOrderLevel    bool
}

var (
	funcMap = template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}
	templates = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/header.html",
		"templates/footer.html", "templates/index.html", "templates/lucky.html",
		"templates/problems.html"))
	serverUrl = "https://s106.net/next2solve"
)

// Render page using a template with data
func renderPage(w http.ResponseWriter, tmpl string, data interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		fmt.Fprintf(w, "Error %v", err)
	}
}

// Show unsolved problems
func showProblems(w http.ResponseWriter, data TemplateData, orderBy string) {
	data.Problems = problems.GetUnsolvedProblems(data.UserID, orderBy)
	if len(data.Problems) == 0 {
		data = TemplateData{UsernameError: "No problems to solve", Username: data.Username}
		renderPage(w, "index", data)
		return
	}
	log.Printf("Get %d problems for user %s\n", len(data.Problems), data.Username)
	renderPage(w, "problems", data)
}

// Show a random unsolved problem
func showRandomProblem(w http.ResponseWriter, data TemplateData) {
	// Choose a problem with lowest dacu, starred first
	data.Problems = problems.GetUnsolvedProblemRandom(data.UserID)
	if len(data.Problems) == 0 {
		data = TemplateData{UsernameError: "No problems to solve", Username: data.Username}
		renderPage(w, "index", data)
		return
	}
	renderPage(w, "lucky", data)
}

// Set template data
func setTemplateData(username string) (TemplateData, error) {
	// Check if username is valid
	userid, err := problems.GetUserID(username)
	if err != nil {
		return TemplateData{UsernameError: err.Error(), Username: username}, err
	}
	// Set user data
	return TemplateData{UserID: userid, Username: username}, nil
}

// Handles requests
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	// Favicon not handled!
	if r.URL.String() == "/favicon.ico" {
		return
	}

	// POST request
	if r.Method == "POST" {
		// Get username
		username := r.PostFormValue("username")
		if r.PostFormValue("show-problems") != "" {
			// Show all unsolved problems
			http.Redirect(w, r, serverUrl+"?all&u="+username+"&o=star", http.StatusFound)

		} else if r.PostFormValue("feeling-lucky") != "" {
			// Show a random unsolved problem
			http.Redirect(w, r, serverUrl+"?lucky&u="+username, http.StatusFound)

		} else {
			// Option not available...
			http.NotFound(w, r)
		}
		return
	}

	// GET request
	query := r.URL.Query()

	// Get username
	if username, ok := query["u"]; ok && len(username[0]) > 0 {
		data, err := setTemplateData(username[0])
		if err != nil {
			data = TemplateData{UsernameError: err.Error(), Username: username[0]}
			renderPage(w, "index", data)
			return
		}

		if _, ok := query["all"]; ok {
			orderBy := ""
			if o, ok := query["o"]; ok && len(o) > 0 {
				orderBy = query["o"][0]
			}
			switch orderBy {
			case "star":
				data.IsOrderStar = true
			case "cat":
				data.IsOrderCategory = true
			case "lev":
				data.IsOrderLevel = true
			}

			// Show all unsolved problems
			showProblems(w, data, orderBy)
			return
		}
		if _, ok := query["lucky"]; ok {
			// Show a random unsolved problem
			showRandomProblem(w, data)
			return
		}
	}

	// GET - Default
	renderPage(w, "index", TemplateData{})
}

// Set handlers, initialize API server and start HTTP server
func HttpServerStart(addr string, apiUrl string) {
	problems.InitAPIServer(apiUrl)
	http.HandleFunc("/", RequestHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
