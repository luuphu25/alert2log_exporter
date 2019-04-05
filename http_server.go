package main

import (
	"io"
	"net/http"
)
import "fmt"
import "encoding/json"
import "time"
import "gopkg.in/olivere/elastic.v6"
import "context"

type test_struct struct {
	Test  string `json:"test"`
	Time  string `json:"timestamp"`
	Value string `json:"value"`
}

func output(t test_struct) {
	fmt.Printf(t.Time + "\n")
	fmt.Printf(t.Test + "\n")
	fmt.Printf(t.Value + "\n")
}
func connectEs(url string, t test_struct) {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		panic(err)
	}
	// Ping to elastic
	info, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	exists, err := client.IndexExists("april5").Do(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		_, err = client.CreateIndex("april5").Do(ctx)
		if err != nil {
			panic(err)
		}
	}
	_, err = client.Index().Index("april5").Type("doc").BodyJson(t).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insert to Elastic success\n")

}
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
func webhook(res http.ResponseWriter, req *http.Request) {
	var url string = "http://192.168.146.131:9200"
	decode := json.NewDecoder(req.Body)
	var t test_struct
	err := decode.Decode(&t)
	if err != nil {
		panic(err)
	}
	t.Time = time.Now().Format(time.RFC3339)
	output(t)
	connectEs(url, t)
	temp := test_struct{Test: "olivere", Time: "123456", Value: "Take Five"}
	output(temp)
	connectEs(url, temp)

}
func main() {
	http.HandleFunc("/metrics", Hello)
	http.HandleFunc("/webhook", webhook)

	fmt.Printf("Server is running at 0.0.0.0:9000\n")
	http.ListenAndServe("127.0.0.1:9000", nil)
}
