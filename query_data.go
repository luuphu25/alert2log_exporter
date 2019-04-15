package main

import (
	"fmt"
	"time"
	"encoding/json"
	"net/http"

)

type test_struct struct {
	Status string `json:"status"`
	Data struct {
		ResultType string `json:"resultType"`
		Result [] struct{
			Metric map[string]string `json:"Metric"`
			Values struct {
				Timestap string `json:"timestamp"`
				Value string `json:"Value`
			} `json:"Values"`
		} `json:"Result"`
	}
}



func main() {
	var time_start_s string
	var time_end_s string
	loc := time.FixedZone("UTC-0", 0)
	//m, _ := time.ParseDuration("5m")
	var metric string = "up"
	var step string = "15s"
	time_end := time.Now().In(loc)
	time_start := time_end.Add(-time.Minute*5)
	time_end_s = time_end.Format(time.RFC3339)
	time_start_s = time_start.Format(time.RFC3339)
	var query string = "http://127.0.0.1:9090/api/v1/query_range?query=" + metric + "&start=" + time_start_s + "&end=" +time_end_s + "&step="+step 
	fmt.Printf(query + "\n")

	req, err := http.Get(query)

	if err != nil {
        panic(err)
    }
	defer req.Body.Close()
	var target test_struct
	json.NewDecoder(req.Body).Decode(&target)
	fmt.Printf(target.Data.Result.Metric)
}