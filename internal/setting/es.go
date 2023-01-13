package setting

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"mall/internal/global"
	"mall/internal/model"
	"os"
)

type es struct {
}

func (es) InitEs() {
	//初始化连接
	host := fmt.Sprintf("http://%s:%d", global.Setting.EsInfo.Host, global.Setting.EsInfo.Port)
	logger := log.New(os.Stdout, "mall", log.LstdFlags)
	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	//新建mapping和index
	exists, err := global.EsClient.IndexExists(model.EsGoods{}.IndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists { // 不存在的时候才需要新建mapping
		_, err = global.EsClient.CreateIndex(model.EsGoods{}.IndexName()).BodyString(model.EsGoods{}.MapJson()).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
