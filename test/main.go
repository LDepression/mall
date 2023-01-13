package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

func main() {
	//初始化连接
	host := fmt.Sprintf("http://%s:%d", "192.168.28.100", 9200)
	logger := log.New(os.Stdout, "mall", log.LstdFlags)
	var err error
	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}
	q := elastic.NewMatchQuery("name", "给过我问问")
	result, _ := client.Search().Index("good").Query(q).Do(context.Background())
	fmt.Println(result.Hits.TotalHits.Value)
}
