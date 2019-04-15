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
type Time int64

type SampleValue float64

type SamplePair struct {
	Timestamp Time       `json:"timestamp"`
	Value     SampleValue `json:"value"`
	
}


func (t *Time) UnmarshalJSON(b []byte) error {
	var minimumTick = time.Millisecond
	var second = int64(time.Second / minimumTick)
	var dotPrecision = int(math.Log10(float64(second)))
	
	
	p := strings.Split(string(b), ".")
	switch len(p) {
	case 1:
		v, err := strconv.ParseInt(string(p[0]), 10, 64)
		if err != nil {
			return err
		}
		*t = Time(v * second)

	case 2:
		v, err := strconv.ParseInt(string(p[0]), 10, 64)
		if err != nil {
			return err
		}
		v *= second

		prec := dotPrecision - len(p[1])
		if prec < 0 {
			p[1] = p[1][:dotPrecision]
		} else if prec > 0 {
			p[1] = p[1] + strings.Repeat("0", prec)
		}

		va, err := strconv.ParseInt(p[1], 10, 32)
		if err != nil {
			return err
		}

		*t = Time(v + va)

	default:
		return fmt.Errorf("invalid time %q", string(b))
	}
	return nil
}


func (v *SampleValue) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("sample value must be a quoted string")
	}
	f, err := strconv.ParseFloat(string(b[1:len(b)-1]), 64)
	if err != nil {
		return err
	}
	*v = SampleValue(f)
	return nil
}
func (s *SamplePair) UnmarshalJSON(b []byte) error {
	v := [...]json.Unmarshaler{&s.Timestamp, &s.Value}
	return json.Unmarshal(b, &v)
}

type test_struct struct {
	Status string `json:"status"`
	Data struct {
		ResultType string `json:"resultType"`
		Result [] struct{
			Metric map[string]string `json:"Metric"`
			Values []SamplePair `json:"values"`
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
	body, err := ioutil.ReadAll(req.Body)
	var target test_struct
	//json.NewDecoder(req.Body).Decode(&target)
	json.Unmarshal(body, &target)
	var url string = "http://127.0.0.1:9200"
	var indexName string = "www"
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