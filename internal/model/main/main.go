package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"io"
	"log"
	"mall/internal/global"
	"mall/internal/model"
	"os"
	"time"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, err := io.WriteString(Md5, code)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(Md5.Sum(nil))
}
func main() {
	//dsn := "root:zxz123456@tcp(127.0.0.1:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags),
	//	logger.Config{
	//		SlowThreshold: time.Second, //慢阙值
	//		Colorful:      true,        //禁用彩色
	//		LogLevel:      logger.Info,
	//	})
	////全局模式
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: newLogger,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//db.AutoMigrate(&model.Category{}, &model.Brands{}, &model.Goods{}, &model.GoodsCategoryBrand{})
	Mysql2Es()

}

//将mysql的数据同步到es中去
func Mysql2Es() {
	var err error
	dsn := "root:zxz123456@tcp(127.0.0.1:3306)/mall?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢阙值
			Colorful:      true,        //禁用彩色
			LogLevel:      logger.Info,
		})
	//全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	//初始化连接
	host := "http://192.168.28.3:9200"
	logger1 := log.New(os.Stdout, "mall", log.LstdFlags)

	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetTraceLog(logger1))
	if err != nil {
		panic(err)
	}

	var goods []model.Good
	db.Find(&goods)
	for _, g := range goods {
		esModel := model.EsGoods{
			ID:          g.ID,
			CategoryID:  g.CategoryID,
			BrandsID:    g.BrandID,
			OnSale:      g.OnSale,
			ShipFree:    g.ShipFree,
			IsNew:       g.IsNew,
			IsHot:       g.IsHot,
			Name:        g.Name,
			ClickNum:    g.ClickNum,
			SoldNum:     g.SoldNum,
			FavNum:      g.FavNum,
			MarketPrice: g.MarketPrice,
			GoodsBrief:  g.GoodsBrief,
			ShopPrice:   g.ShopPrice,
		}
		//第一个index表示的是名词,第二index表示的动词,后面的参数就是表的名字,后面加上id的话,这样就会在我们es中加上
		//加入一条数据的固定写法
		_, err2 := global.EsClient.Index().Index(model.EsGoods{}.IndexName()).Id(utils.ToString(g.ID)).BodyJson(esModel).Do(context.Background())
		if err2 != nil {
			panic(err)
		}
	}
}
