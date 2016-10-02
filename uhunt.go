// Next problem to solve
// https://github.com/maru/next2solve
//
// uHunt API calls

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
  "errors"
)

// URL paths of uHunt API
// First argument is the server address, e.g. "http://uhunt.felix-halim.net"
const (
	APIUsernameToUserid = "%s/api/uname2uid/%s"
	APIUserSubmissions = "%s/api/subs-user/%s"
	APIUserSubmissionsToProblems = "%s/api/subs-pids/%s/%s/0"
	APIProblemList = "%s/api/p"
)

//
//  Get userid by username, output error if username is not found.
//
func apiGetUserID(w http.ResponseWriter, username string) (string, error) {
	url := fmt.Sprintf(APIUsernameToUserid, APIUrl, username)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return "", err
	}
	id := string(body)
	if id == "0" {
		return "", errors.New("Username not found")
	}
	return id, nil
}

// func Api
// <!-- "dacu":7895,
// "sube":153,"noj":0,"inq":0,"ce":1055,"rf":0,"re":774,"ole":9,"tle":1535,"mle":1,"wa":7639,"pe":17,"ac":10377,
// 153+1055+774+1535+9+1+7639+10377+17 -->
