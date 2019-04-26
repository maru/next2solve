// Next problem to solve
// https://github.com/maru/next2solve
//
// Shows unsolved problems from UVa Online Judge, using the uHunt API.
//
// Main routine
//
// Usage of ./next2solve:
//   -api string
//     	API URL (default "https://uhunt.onlinejudge.org")
//   -p string
//     	Listening port (default "8002")

package main

import (
	"flag"
	"next2solve/server"
)

func main() {
	APIUrl := flag.String("api", "https://uhunt.onlinejudge.org", "API URL")
	port := flag.String("p", "8002", "Listening port")
	flag.Parse()
	*port = "127.0.0.1:" + *port
	server.HttpServerStart(*port, *APIUrl)
}
