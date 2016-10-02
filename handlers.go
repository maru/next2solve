// Next problem to solve
// https://github.com/maru/next2solve
//
// HTTP handlers

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type UserInfo struct {
	UsernameError string
	UserID        string
	Username      string
}

var (
	APIUrl string
	templates = template.Must(template.ParseFiles("html/header.html",
		"html/footer.html", "html/index.html", "html/lucky.html", "html/problems.html"))
)

//
// Render page using a template with data
//
func renderPage(w http.ResponseWriter, tmpl string, data interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		fmt.Fprintf(w, "Error %v", err)
	}
}

//
// Get user information from cookies
//
func getUserInfo(r *http.Request) (UserInfo) {
	userInfo := UserInfo{}
	if cookie, err := r.Cookie("userid"); err == nil {
		userInfo.UserID = cookie.Value
	}
	if cookie, err := r.Cookie("username"); err == nil {
		userInfo.Username = cookie.Value
	}
	return userInfo
}

//
// Set user information in cookies
//
func setUserInfo(w http.ResponseWriter, userInfo *UserInfo, userid, username string) {
	if userInfo.UserID != userid {
		cookie := http.Cookie{Name: "userid", Value: userid}
		http.SetCookie(w, &cookie)
		userInfo.UserID = userid
	}
	if userInfo.Username != username {
		cookie := http.Cookie{Name: "username", Value: username}
		http.SetCookie(w, &cookie)
		userInfo.Username = username
	}
}

//
// Handles requests
//
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	userInfo := getUserInfo(r)
	if r.Method == "POST" {
    // Show problems to solve
		username := r.PostFormValue("username")
    // Check if username is valid
    userid, err := apiGetUserID(w, username)
		if err != nil {
			renderPage(w, "index", UserInfo{err.Error(), "", username})
      return
		}
		// Set user information in a cookie
		setUserInfo(w, &userInfo, userid, username)

		// Show all unsolved problems
		if r.PostFormValue("show-problems") != "" {
			// Show all unsolved problems
			showProblems(w, userInfo)
			return
		} else {
			// Show a random unsolved problem
			showRandomProblem(w, userInfo)
			return
		}
	}

	// GET - Default
	renderPage(w, "index", userInfo)
}

//
// Set handlers and start http server
//
func httpServerStart(addr string, apiUrl string) {
	APIUrl = apiUrl
	http.HandleFunc("/", RequestHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
