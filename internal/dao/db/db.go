package db

import (
	"fmt"
	"log"
	"mall/internal/global"
	"mall/internal/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Init() *gorm.DB {
	m := global.Setting.Mysql
	fmt.Println(m)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, m.DbName)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢阙值
			Colorful:      true,        //禁用彩色
			LogLevel:      logger.Info,
		})
	//全局模式
	var err error
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})
	if err != nil {
		panic(err)
	}
	_ = DB.AutoMigrate(&model.RebackDetails{})
	//for i := 421; i <= 600; i++ {
	//	var inv model.Inventory
	//	inv.Goods = int32(i)
	//	inv.Stocks = 100
	//	inv.Version = 0
	//	DB.Create(&inv)
	//}
	/*
		DB.Create(&model.RebackDetails{
			GoodsInfo: []model.GoodInfo{
				{GoodID: 421, Num: 2},
				{GoodID: 422, Num: 2},
			},
			OrderSn: "lyc",
			Status:  1,
		})

	*/

	/*
		var reback model.RebackDetails
		DB.Where(&model.RebackDetails{OrderSn: "lyc"}).First(&reback)
		fmt.Println(reback)

	*/
	return DB
}
