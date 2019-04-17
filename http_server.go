package main

import (
	//"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"os"
	"io"
	"net/http"
	//"time"
	"github.com/olivere/elastic"
	"github.com/luuphu25/alert2log_exporter/model"
	"github.com/luuphu25/alert2log_exporter/query"
	"github.com/luuphu25/alert2log_exporter/lib"
)
// Hello path testing 
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
	var metric string = "node_memory_MemFree_bytes"
	var step string = "15s"
	prometheus_url := "http://127.0.0.1:9090"
	var alert_receive model.Notification //alert struct
	var past_data  model.Query_struct // past data struct
	var complete_data model.Log_Data // merge struct
	decode := json.NewDecoder(req.Body)
	
	// create variable (type struct Notification) to handle data send by alertmanage
		err := decode.Decode(&alert_receive)
	if err != nil {
		panic(err)
	}
	
	// Query data from prometheus	
	start_time, end_time := query.CreateTime(alert_receive.Alerts[0].StartsAt)
	fmt.Printf(start_time + " : " + end_time)
	var query_command = query.CreateQuery(prometheus_url, metric, start_time, end_time, step)
	//fmt.Printf("Query: " + query_command + "\n")
	past_data = query.Query_past(query_command)

	// merge alert data + query data 
	complete_data.AlertInfo = alert_receive
	complete_data.PastData = past_data

	// insert into Elastich
	lib.InsertEs(client, complete_data, indexName)

	// write log file
	file, _ := json.MarshalIndent(complete_data, "", " ")
	lib.WriteFile(file)
	

}

// print data post from prometheus

func getAlert(client *elastic.Client, indexName string, req *http.Request){
	//resp, err := http.Get("http://127.0.0.1:9093/api/v2/alerts")
	fmt.Printf("Get alert pending \n")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil{
		panic(err)
	}
	bodyString := string(body)
	fmt.Printf(bodyString + "\n")

}
func main() {
	var url string = "http://127.0.0.1:9200"
	var indexName string = "log_db"
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
