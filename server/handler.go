// Next problem to solve
// https://github.com/maru/next2solve
//
// HTTP handler

package server

import (
	"log"
	"net"
	"net/http"
	"next2solve/problems"
	"os"
)

// Handles requests
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	log.Printf("[%v] %s %s\n", ip, r.Method, r.URL.Path)

	switch r.URL.Path {
	case "/":
		IndexHandler(w, r)
	case "/all":
		AllHandler(w, r)
	case "/lucky":
		LuckyHandler(w, r)
	case "/favicon.ico":
		http.ServeFile(w, r, "templates/favicon.ico")
	default:
		http.NotFound(w, r)
	}
}

// Show unsolved problems
func AllHandler(w http.ResponseWriter, r *http.Request) {
	var data TemplateData
	var username []string
	query := r.URL.Query()

	// Check username field
	var ok bool
	if username, ok = query["u"]; !ok || len(username[0]) == 0 {
		data = TemplateData{UsernameError: "Please enter your UVa username"}
		renderPage(w, "index", data)
		return
	}

	// Validate user and set data to template
	var err error
	if data, err = setTemplateUserData(username[0]); err != nil {
		renderPage(w, "index", data)
		return
	}

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
	data.Problems = problems.GetUnsolvedProblems(data.UserID, orderBy)
	if len(data.Problems) == 0 {
		data = TemplateData{UsernameError: "No problems to solve!"}
		renderPage(w, "index", data)
		return
	}

	log.Printf("Get %d problems for user %s\n", len(data.Problems), data.Username)
	renderPage(w, "problems", data)
}

// Show a random unsolved problem
// Choose a problem with lowest dacu, starred first
func LuckyHandler(w http.ResponseWriter, r *http.Request) {
	var data TemplateData
	var username []string
	query := r.URL.Query()

	// Check username field
	var ok bool
	if username, ok = query["u"]; !ok || len(username[0]) == 0 {
		data = TemplateData{UsernameError: "Please enter your UVa username"}
		renderPage(w, "index", data)
		return
	}

	// Validate user and set data to template
	var err error
	if data, err = setTemplateUserData(username[0]); err != nil {
		renderPage(w, "index", data)
		return
	}

	// Show a random unsolved problem
	data.Problems = problems.GetUnsolvedProblemRandom(data.UserID)
	if len(data.Problems) == 0 {
		data = TemplateData{UsernameError: "No problems to solve!"}
		renderPage(w, "index", data)
		return
	}

	renderPage(w, "lucky", data)
}

// Index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// POST request
	if r.Method == "POST" {
		// Get username
		username := r.PostFormValue("username")
		if r.PostFormValue("show-problems") != "" {
			// Show all unsolved problems
			http.Redirect(w, r, "/all?u="+username+"&o=star", http.StatusFound)

		} else if r.PostFormValue("feeling-lucky") != "" {
			// Show a random unsolved problem
			http.Redirect(w, r, "/lucky?u="+username, http.StatusFound)
		}
	} else {
		// GET - Default
		renderPage(w, "index", TemplateData{})
	}
}

// Set handlers, initialize API server and start HTTP server
func HttpServerStart(addr string, apiUrl string, logfile string) {
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatal("Could not create logfile ", logfile)
		}
		log.SetOutput(file)
	}
	LoadTemplates()
	problems.InitAPIServer(apiUrl)
	http.HandleFunc("/", ServeHTTP)
	log.Fatal(http.ListenAndServe(addr, nil))
}
