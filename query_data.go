package main

import (
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"github.com/olivere/elastic"
	"context"

)
type temp struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name     string `json:"__name__"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}
type SampleValue float64

type SamplePair struct {
	Timestamp time.Time       `json:"timestamp"`
	Value     SampleValue `json:"value"`
	
}
type test_struct struct {
	Status string `json:"status"`
	Data struct {
		ResultType string `json:"resultType"`
		Result [] struct{
			Metric map[string]string `json:"Metric"`
			Values [] SamplePair `json:"values"`
		} `json:"Result"`
	}
}



func main() {
	var time_start_s string
	var time_end_s string
	loc := time.FixedZone("UTC-0", 0)
	//m, _ := time.ParseDuration("5m")
	var metric string = "node_memory_MemFree_bytes"
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

	var url string = "http://127.0.0.1:9200"
	var indexName string = "testdb"
	//create client elastic
	client, err := elastic.NewClient(elastic.SetURL(url))

	if err != nil{
		panic(err)
	}
	InsertEs(client, target, indexName)
}

func InsertEs(client *elastic.Client, data test_struct, indexName string){
	ctx := context.Background()
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		_, err = client.CreateIndex(indexName).Do(ctx)
		if err != nil {
			panic(err)
		}
	}
	_, err = client.Index().Index(indexName).Type("doc").BodyJson(data).Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nInsert to Elastic success\n")
}