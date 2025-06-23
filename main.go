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

type APIResponse struct {
	Status string    `json:"status"`
	Result []Contest `json:"result"`
}

func main() {
	url := "https://codeforces.com/api/contest.list?gym=false"
	res, err := make_request(url)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
		return
	}
	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	// fmt.Print(*res)
	print_contests(apiResp)

}

func make_request(url string) (*http.Response, error) {
	return http.Get(url)
}

func print_contests(apiResp APIResponse) {
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
