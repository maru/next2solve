package main

import (
        "io"
        "io/ioutil"
        "fmt"
        "net/http"
)

// Load file
func loadPage(filename string) ([]byte, error) {
    page, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return page, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    io.WriteString(w, "Hello world!")
    return
  }

  // Default index.html
  p, err := loadPage("index.html")
  if err != nil {
    fmt.Fprintf(w, "Error %v", err)
    return
  }
  fmt.Fprintf(w, string(p))
}

func main() {
  fmt.Println("Web server ON")
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8002", nil)
}
