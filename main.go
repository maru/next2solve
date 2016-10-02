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
//     	API URL (default "http://uhunt.felix-halim.net")

package main

import "flag"

func main() {
	APIUrl := flag.String("api", "http://uhunt.felix-halim.net", "API URL")
	addr := flag.String("addr", ":8002", "listening address")
	flag.Parse()
	httpServerStart(*addr, *APIUrl)
}
