// Next problem to solve
// https://github.com/maru/next2solve
//
// Shows unsolved problems from UVa Online Judge, using the uHunt API.
//
// Main routine
//
// Usage of ./next2solve:
//   -addr string
//     	listening address (default ":8002")
//   -api string
//     	API URL (default "https://uhunt.onlinejudge.org")

package main

import "flag"

func main() {
	APIUrl := flag.String("api", "https://uhunt.onlinejudge.org", "API URL")
	addr := flag.String("addr", ":8002", "listening address")
	flag.Parse()
	HttpServerStart(*addr, *APIUrl)
}
