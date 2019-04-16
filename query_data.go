package main

import (
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"github.com/olivere/elastic"
	"context"
	"io/ioutil"
	"github.com/luuphu25/alert2log_exporter/model"



)

func main() {
	loc := time.FixedZone("UTC-0", 0)
	//m, _ := time.ParseDuration("5m")
	var metric string = "mem_used"
	var step string = "15s"
	time_end := time.Now().In(loc)
	time_start := time_end.Add(-time.Minute*5) //for 5m
	time_end_s := time_end.Format(time.RFC3339)
	time_start_s := time_start.Format(time.RFC3339)
	url := "http://61.28.251.119:9090"
	query := CreateQuery(url, metric, time_start_s, time_end_s, step)
	
	fmt.Printf(query + "\n")
	req, err := http.Get(query)


	if err != nil {
        panic(err)
    }
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	var target model.Query_struct
	//json.NewDecoder(req.Body).Decode(&target)
	json.Unmarshal(body, &target)
	fmt.Printf(target.Data.Result[0].Values[0].String()) // test print SamplePair

	var es_url string = "http://127.0.0.1:9200"
	var indexName string = "test_query"
	//create client elastic
	client, err := elastic.NewClient(elastic.SetURL(es_url))

	if err != nil{
		panic(err)
	} 
	InsertEs(client, target, indexName)
}

func CreateQuery(url string, metric string, time_start_s string, time_end_s string, step string) string {
	var query string
	query  = url + "/api/v1/query_range?query=" + metric + "&start=" + time_start_s + "&end=" +time_end_s + "&step="+ step 
	return query
}


func InsertEs(client *elastic.Client, data model.Query_struct, indexName string){
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