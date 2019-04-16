package main

import (
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"github.com/olivere/elastic"
	"context"
	"strconv"
	"io/ioutil"
	"strings"
	"math"
	"github.com/luuphu25/alert2log_exporter/model"



)
func create_query(url string, metric string, time_start_s string, time_end_s string, step string){
	var query string = url + "/api/v1/query_range?query=" + metric + "&start=" + time_start_s + "&end=" +time_end_s + "&step="+step 
	return query

}

func main() {
	//var time_start_s string
	//var time_end_s string
	loc := time.FixedZone("UTC-0", 0)
	//m, _ := time.ParseDuration("5m")
	var metric string = "up"
	var step string = "15s"
	time_end := time.Now().In(loc)
	time_start := time_end.Add(-time.Minute*5)
	time_end_s := time_end.Format(time.RFC3339)
	time_start_s := time_start.Format(time.RFC3339)
	url := "http://127.0.0.1:9090"
	var query string = create_query(metric, time_start_s, time_end_s, step)
	req, err := http.Get(query)

	if err != nil {
        panic(err)
    }
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	var target model.Query_struct
	//json.NewDecoder(req.Body).Decode(&target)
	json.Unmarshal(body, &target)
	//var url string = "http://127.0.0.1:9200"
	//var indexName string = "www"
	//create client elastic
	client, err := elastic.NewClient(elastic.SetURL(url))

	if err != nil{
		panic(err)
	}
	//InsertEs(client, target, indexName)
	fmt.Printf(target.)
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