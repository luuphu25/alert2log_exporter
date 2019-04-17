package lib

import (
	//"io/ioutil"
	"log"
	"os"
	"time"
	"fmt"
	"github.com/olivere/elastic"
	"context"
)

func WriteFile(data []byte) {
	timestamp := int32(time.Now().Unix())
	times := []byte(fmt.Sprintf("%d", timestamp))
	date := time.Now().UTC().Format("01-02-2006")
	var filename = "./Log/logAlert_" + date + ".json"
	_, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err){
			_, err := os.Create(filename)
			if err != nil {
				log.Fatal("Can't create log file", err)
			}
		}
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Can't open new file", err)
	}


	defer f.Close()
	if _, err = f.Write(times); err != nil {
		log.Fatal("Can't write timestamp to file", err)
	}
	if _, err = f.Write(data); err != nil {
		log.Fatal("Can't write to file", err)
	}
	fmt.Printf("Write data to file success!\n")
	
}

func InsertEs(client *elastic.Client, data model.Log_Data, indexName string){
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
