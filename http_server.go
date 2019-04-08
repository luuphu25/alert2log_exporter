package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
	"time"
	"github.com/olivere/elastic"
	"github.com/luuphu25/alert2log_exporter/model"
)
// Insert data in elastic
func InsertEs(client *elastic.Client, data model.Notification, indexName string){
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

// Hello path 
func Hello(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	io.WriteString(
		res,
		`<h1>Hello World</h1>`,
	)
}
// handle email send by alertmanager 
func webhook(client *elastic.Client, indexName string, req *http.Request) {
	
	decode := json.NewDecoder(req.Body)
	// create variable (type struct Notification) to handle data send by alertmanage
	var t model.Notification
	err := decode.Decode(&t)
	if err != nil {
		panic(err)
	}
	t.Timestamp = time.Now().Format(time.RFC3339)
	InsertEs(client, t, indexName)

}

// print data post from prometheus

func getAlert(client *elastic.Client, indexName string, req *http.Request){
	//resp, err := http.Get("http://127.0.0.1:9093/api/v2/alerts")
	fmt.Printf("Get alert from prometheus\n")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil{
		panic(err)
	}
	bodyString := string(body)
	fmt.Printf(bodyString)

}
func main() {
	var url string = "http://192.168.146.131:9200"
	var indexName string = "testdb"
	//create client elastic
	client, err := elastic.NewClient(elastic.SetURL(url))

	if err != nil{
		panic(err)
	}

	// handle request
	http.HandleFunc("/api/v1/alerts", func(res http.ResponseWriter, req *http.Request){
		getAlert(client, "getSignal", req)
	})
	http.HandleFunc("/metrics", Hello)
	http.HandleFunc("/webhook", func(res http.ResponseWriter, req *http.Request){
		webhook(client, indexName, req)
	})

	fmt.Printf("Server is running at 0.0.0.0:9000\n")
	http.ListenAndServe("0.0.0.0:9000", nil)
}
