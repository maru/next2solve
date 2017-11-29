package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UvaUser struct {
	Rank                     int    `json:"rank"`
	Old                      int    `json:"old"`
	Userid                   int    `json:"userid"`
	Name                     string `json:"name"`
	Username                 string `json:"username"`
	NumberOfAcceptedProblems int    `json:"ac"`
	NumberOfSubmissions      int    `json:"nos"`
	Activity                 []int  `json:"activity"`
}

func getResponse(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func main() {
	url := "http://uhunt.onlinejudge.org/api/ranklist/46232/0/0"
	url314 := "http://uhunt.onlinejudge.org/api/rank/314/1"

	resp, err := getResponse(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Parse the data
	var user []UvaUser
	if err = json.Unmarshal(resp, &user); err != nil {
		fmt.Println(err)
		return
	}

	if user[0].Rank == 314 {
		fmt.Println("Chicapi is 314 :)")
		return
	}
	// else: how many problems to be 314?
	resp, err = getResponse(url314)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Parse the data
	var user2 []UvaUser
	if err = json.Unmarshal(resp, &user2); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Need %d more problems!\n", user2[0].NumberOfAcceptedProblems - user[0].NumberOfAcceptedProblems)
	return
}
