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
  APIUsernameToUserid = "http://uhunt.felix-halim.net/api/uname2uid/%s"
)

type UsernameInfo struct {
  UsernameError  string
  UserID string
  Username string
}

/*
 * Return default index.html
 */
func respIndex(w http.ResponseWriter, data UsernameInfo) {
  t, err := template.ParseFiles("index.html")
  if err != nil {
    fmt.Fprintf(w, "Error %v", err)
    return
  }
  t.Execute(w, data)
}

/*
 * Return default index.html
 */
func respLucky(w http.ResponseWriter, data UsernameInfo) {
  t, err := template.ParseFiles("lucky.html")
  if err != nil {
    fmt.Fprintf(w, "Error %v", err)
    return
  }
  t.Execute(w, data)
}

/*
 *  Get userid by username, output error if username is not found.
 */
func getUserID(w http.ResponseWriter, username string) string {
  url := fmt.Sprintf(APIUsernameToUserid, username)
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
    respIndex(w, UsernameInfo{"Username not found", "", username})
    return ""
  }
  return id
}

/*
 * Handles request to index page
 */
func IndexHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    // Show problems to solve
    var id string
    username := r.PostFormValue("username")
    // Check if username is valid
    if id = getUserID(w, username); id == "" {
        return
    }
    cookie := http.Cookie{ Name: "userid", Value: id }
    http.SetCookie(w, &cookie)
    // Show all unsolved problems
    if r.PostFormValue("show-problems") != "" {

    } else {
      // Show a problem
      respLucky(w, UsernameInfo{"", id, username})
      return
    }
    respIndex(w, UsernameInfo{"", id, username})
    return
  }

  // GET - Default
  respIndex(w, UsernameInfo{})
}

/*
 * Set handlers and start http server
 */
func httpServerStart(addr string) {
  http.HandleFunc("/", IndexHandler)
  log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
  httpServerStart(":8002")
}
