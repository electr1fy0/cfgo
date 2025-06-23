package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Contest struct {
	ID               int    `json:"id"`
	Name             string `json:"name`
	StartTimeSeconds int64  `json:"startTimeSeconds`
}

type Rating struct {
}

type User struct{}

type Submission struct{}
type Problem struct{}

type APIResponse[T any] struct {
	Status string `json:"status"`
	Result []T    `json:"result"`
}

func main() {
	args := os.Args

	if len(args) < 1 {
		fmt.Println("Please provide arguments")
		return
	}

	switch args[1] {
	case "--contests":
		PrintContests()
	case "--rating":

	case "--submissions":
	case "info":

	}

}

func Request(url string) (*http.Response, error) {
	return http.Get(url)
}

func PrintContests() {
	url := "https://codeforces.com/api/contest.list?gym=false"
	res, err := Request(url)
	if err != nil {
		fmt.Println("Error making the request: ", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}
	var apiResp APIResponse[Contest]

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println("Error unmarshalling: ", err)

	}

	table := tablewriter.NewWriter((os.Stdout))
	table.Header([]string{"Contest ID", "Name", "Start time"})

	for _, contest := range apiResp.Result[:5] {
		startTime := time.Unix(contest.StartTimeSeconds, 0).Local().Format("02 Jan 2006 15:04")
		row := []string{
			fmt.Sprintf("%d", contest.ID), contest.Name, startTime,
		}
		table.Append(row)

	}
	table.Render()

}
