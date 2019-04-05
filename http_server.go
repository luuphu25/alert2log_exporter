package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gopkg.in/olivere/elastic.v6"
)

type KV map[string]string
type Data struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   Alert  `json:"alerts"`

	GroupLabels       KV `json:"groupLabels"`
	CommonLabels      KV `json:"commonLabels"`
	CommonAnnotations KV `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

// Alert holds one alert for notification templates.
type Alert struct {
	Status      string    `json:"status"`
	Labels      KV        `json:"labels"`
	Annotations KV        `json:"annotations"`
	StartsAt    time.Time `json:"startsAt"`
	EndsAt      time.Time `json:"endsAt"`
}

type test_struct struct {
	Test  string `json:"test"`
	Time  string `json:"timestamp"`
	Value string `json:"value"`
}
type Alerts []Alert

type notification struct {
	Alerts []struct {
		Annotations  map[string]string `json:"annotations"`
		EndsAt       time.Time         `json:"endsAt"`
		GeneratorURL string            `json:"generatorURL"`
		Labels       map[string]string `json:"labels"`
		StartsAt     time.Time         `json:"startsAt"`
		Status       string            `json:"status"`
	} `json:"alerts"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	CommonLabels      map[string]string `json:"commonLabels"`
	ExternalURL       string            `json:"externalURL"`
	GroupLabels       map[string]string `json:"groupLabels"`
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`

	// Timestamp records when the alert notification was received
	Timestamp string `json:"@timestamp"`
}

func output(t test_struct) {
	fmt.Printf(t.Time + "\n")
	fmt.Printf(t.Test + "\n")
	fmt.Printf(t.Value + "\n")
}
func connectEs(url string, t notification) {
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

	exists, err := client.IndexExists("testing").Do(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		_, err = client.CreateIndex("testing").Do(ctx)
		if err != nil {
			panic(err)
		}
	}
	_, err = client.Index().Index("testing").Type("doc").BodyJson(t).Do(ctx)
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
	var t notification
	err := decode.Decode(&t)
	if err != nil {
		panic(err)
	}
	t.Timestamp = time.Now().Format(time.RFC3339)
	//output(t)
	connectEs(url, t)

}
func main() {
	http.HandleFunc("/metrics", Hello)
	http.HandleFunc("/webhook", webhook)

	fmt.Printf("Server is running at 0.0.0.0:9000\n")
	http.ListenAndServe("0.0.0.0:9000", nil)
}
