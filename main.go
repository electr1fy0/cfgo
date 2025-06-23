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
	Name             string `json:"name"`
	StartTimeSeconds int64  `json:"startTimeSeconds`
}

type Rating struct {
	ContestID   int    `json:"contestID"`
	ContestName string `json:"contestName"`
	Rank        int    `json:"rank"`
	Handle      string `json:"handle"`
	OldRating   int    `json:"oldRating"`
	NewRating   int    `json:"newRating"`
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
		PrintRatingHistory(args[2])
	case "--submissions":
	case "info":

	}

}

func Request(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making the request: ", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		os.Exit(1)
	}
	return body

}

func PrintContests() {
	body := Request("https://codeforces.com/api/contest.list?gym=false")

	var apiResp APIResponse[Contest]

	err := json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println("Error unmarshalling: ", err)
	}

	table := tablewriter.NewWriter((os.Stdout))
	table.Header([]string{"Contest ID", "Name", "Start time"})

	for _, contest := range apiResp.Result[:min(10, len(apiResp.Result))] {
		startTime := time.Unix(contest.StartTimeSeconds, 0).Local().Format("02 Jan 2006 15:04")
		row := []string{
			fmt.Sprintf("%d", contest.ID), contest.Name, startTime,
		}
		table.Append(row)

	}
	table.Render()
}

func PrintRatingHistory(handle string) {
	url := fmt.Sprintf("https://codeforces.com/api/user.rating?handle=%s", handle)
	body := Request(url)

	var apiResp APIResponse[Rating]
	err := json.Unmarshal(body, &apiResp)

	if err != nil {
		fmt.Println("Error umarshalling: ", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{
		"Contest ID",
		"Title",
		// "Handle",
		"Rank",
		"Old Rating",
		"New Rating"})

	limit := len(apiResp.Result)

	for i := limit - 1; i >= 0; i-- {
		ratingItem := apiResp.Result[i]
		var row []string = []string{
			fmt.Sprintf("%d", ratingItem.ContestID),
			ratingItem.ContestName,
			// ratingItem.Handle,
			fmt.Sprintf("%d", ratingItem.Rank),
			fmt.Sprintf("%d", ratingItem.OldRating),
			fmt.Sprintf("%d", ratingItem.NewRating),
		}

		table.Append(row)
	}
	table.Render()

}
