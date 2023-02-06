package main

//func main() {
//	//初始化连接
//	host := fmt.Sprintf("http://%s:%d", "192.168.28.6", 9200)
//	logger := log.New(os.Stdout, "mall", log.LstdFlags)
//	var err error
//	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
//		elastic.SetTraceLog(logger))
//	if err != nil {
//		panic(err)
//	}
//	q := elastic.NewMatchQuery("name", "给过我问问")
//	result, _ := client.Search().Index("good").Query(q).Do(context.Background())
//	fmt.Println(result.Hits.TotalHits.Value)
//}

//func main() {
//	for i := 421; i <= 600; i++ {
//		var inv model.Inventory
//		inv.Goods = int32(i)
//		inv.Stocks = 100
//		inv.Version = 0
//		dao.Group.DB.Create(&inv)
//	}
//}
