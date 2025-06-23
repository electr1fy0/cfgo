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

type APIResponse[T any] struct {
	Status string `json:"status"`
	Result []T    `json:"result"`
}

type Contest struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	StartTimeSeconds int64  `json:"startTimeSeconds`
}

type RatingHistory struct {
	ContestID   int    `json:"contestID"`
	ContestName string `json:"contestName"`
	Rank        int    `json:"rank"`
	Handle      string `json:"handle"`
	OldRating   int    `json:"oldRating"`
	NewRating   int    `json:"newRating"`
}

type User struct {
	Rank      string `json:"rank"`
	Handle    string `json:"handle"`
	MaxRating int    `json:"maxRating"`
	Rating    int    `json:"rating"`
}

// Following two structs work together
type Submission struct {
	ContestID           int     `json:"contestId"`
	CreationTimeSeconds int64   `json:"creationTimeSeconds"`
	Problem             Problem `json:"problem"`
	Verdict             string  `json:"verdict"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
}

type Problem struct {
	Name   string `json:"name"`
	Index  string `json:"index"`
	Rating *int   `json:"rating"` // optional (nil if missing)
}

func PrintUsage() {
	fmt.Print(`
Usage: cfetch <--rating <RATING>|--info <INFO>|--contests|--submissions <SUBMISSIONS>>

Options:
  -r, --rating <RATING>
  -i, --info <INFO>
  -c, --contests
  -s, --submissions <SUBMISSIONS> `)
	os.Exit(1)
}
func main() {
	args := os.Args

	if len(args) < 2 {
		PrintUsage()
	}

	flag := args[1]

	if flag != "--contests" && flag != "-c" && len(args) < 3 {
		PrintUsage()
	}

	switch args[1] {
	case "--contests", "-c":
		PrintContests()
	case "--rating", "-r":

		PrintRatingHistory(args[2])
	case "--submissions", "-s":
		PrintSubmissionHistory(args[2])
	case "--info", "-i":
		PrintUserInfo(args[2])
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

	var apiResp APIResponse[RatingHistory]
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

func PrintUserInfo(handle string) {
	url := fmt.Sprintf("https://codeforces.com/api/user.info?handles=%s&checkHistoricHandles=false", handle)
	body := Request(url)

	var apiResp APIResponse[User]
	err := json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Handle", "Rank", "Rating", "Max Rating"})

	for _, user := range apiResp.Result {
		row := []string{
			user.Handle,
			user.Rank,
			fmt.Sprintf("%d", user.Rating),
			fmt.Sprintf("%d", user.MaxRating),
		}
		table.Append(row)
	}

	table.Render()
}

func PrintSubmissionHistory(handle string) {
	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=10", handle)

	body := Request(url)

	var apiResp APIResponse[Submission]

	err := json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println("Error unmarshalling: ", err)

	}

	table := tablewriter.NewTable(os.Stdout)
	table.Header([]string{
		"Contest ID",
		"Difficulty",
		"Problem Name",
		"Verdict",
		"Language",
		"Time",
	})

	for _, submission := range apiResp.Result {
		var difficulty string
		if submission.Problem.Rating == nil {
			difficulty = "N/A"
		} else {
			difficulty = fmt.Sprintf("%d", submission.Problem.Rating)
		}

		startTime := time.Unix(submission.CreationTimeSeconds, 0).Local().Format("02 Jan 2006 15:04")

		var row []string = []string{
			fmt.Sprintf("%d", submission.ContestID),
			difficulty,
			submission.Problem.Name,
			submission.Verdict,
			submission.ProgrammingLanguage,
			startTime,
		}
		table.Append(row)
	}
	table.Render()
}
