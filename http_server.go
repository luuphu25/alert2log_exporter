package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/olivere/elastic"
	"alert2log_exporter/model"
)

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
	fmt.Printf("Insert to Elastic success\n")
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
// webhook path 
func webhook(client *elastic.Client, indexName string, req *http.Request) {
	
	decode := json.NewDecoder(req.Body)
	var t model.Notification
	err := decode.Decode(&t)
	if err != nil {
		panic(err)
	}
	t.Timestamp = time.Now().Format(time.RFC3339)
	//output(t)
	InsertEs(client, t, "alert_test")

}
func main() {
	var url string = "http://192.168.146.131:9200"
	var indexName string = "alert_1"
	client, err := elastic.NewClient(elastic.SetURL(url))

	if err != nil{
		panic(err)
	}

	// handle request
	http.HandleFunc("/metrics", Hello)
	http.HandleFunc("/webhook", func(res http.ResponseWriter, req *http.Request){
		webhook(client, indexName, req)
	})

	fmt.Printf("Server is running at 0.0.0.0:9000\n")
	http.ListenAndServe("0.0.0.0:9000", nil)
}
