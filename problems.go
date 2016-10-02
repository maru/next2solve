// Next problem to solve
// https://github.com/maru/next2solve
//
// Problems
//
// Number of Distinct Accepted User (DACU)

package main

import (
	// "fmt"
	// "io/ioutil"
	"net/http"
)

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
