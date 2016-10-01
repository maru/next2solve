// dacu: different accepted users

package main

import (
	// "bytes"
	// "encoding/binary"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	APIUsernameToUserid = "%s/api/uname2uid/%s"
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
//  Get userid by username, output error if username is not found.
//
func getUserID(w http.ResponseWriter, username string) string {
	url := fmt.Sprintf(APIUsernameToUserid, APIUrl, username)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return ""
	}
	id := string(body)
	if id == "0" {
		renderPage(w, "index", UserInfo{"Username not found", "", username})
		return ""
	}
	return id
}

//
//
//
func showProblems(w http.ResponseWriter, userInfo UserInfo) {
	 //Math.max(1, 10 - Math.floor(Math.min(10, Math.log(p.dacu))));
	 renderPage(w, "problems", userInfo)
}

//
//
//
func showRandomProblem(w http.ResponseWriter, userInfo UserInfo) {
	 // Choose a problem with lowest dacu, starred first

	 renderPage(w, "lucky", userInfo)
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
		var userid string
		username := r.PostFormValue("username")
		// Check if username is valid
		if userid = getUserID(w, username); userid == "" {
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
func httpServerStart(addr string) {
	APIUrl = "http://uhunt.felix-halim.net"
	http.HandleFunc("/", RequestHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
