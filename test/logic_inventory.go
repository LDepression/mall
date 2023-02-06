package main

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mall/internal/form"
	"mall/internal/model"

	"mall/internal/pkg/app/errcode"
	"os"
	"sync"
	"time"
)

func Sell(wg *sync.WaitGroup, db *gorm.DB, req form.SellInfo) errcode.Err {
	//client := goredislib.NewClient(&goredislib.Options{
	//	Addr:     "localhost:6379",
	//	Password: "123456",
	//})
	//pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	//rs := redsync.New(pool)
	tx := db.Begin()
	defer wg.Done()
	for _, GoodInfo := range req.GoodsInfo {
		var inv model.Inventory
		for {
			//mutex := rs.NewMutex(fmt.Sprintf("goods_%d", GoodInfo.GoodID)) //不同的商品
			//if err := mutex.Lock(); err != nil {
			//	return errcode.ErrServer.WithDetails("释放redis分布式锁异常")
			//}
			//先去查询一下
			if result := db.Where("goods=?", GoodInfo.GoodID).Find(&inv); result.RowsAffected == 0 {
				//说明此时没哟这个消息的记录,那么此时扣减库存就失败了
				tx.Rollback()
				return errcode.ErrServer.WithDetails(result.Error.Error())
			}
			fmt.Println(inv)
			//这里先去判断一下,商品的库存是否够用
			//这里表示的是如果此时库存不够
			if inv.Stocks < GoodInfo.Num {
				tx.Rollback()
				return errcode.ErrServer.WithDetails("商品的库存不足....")
			}
			inv.Stocks -= GoodInfo.Num //扣减库存
			//update inventory set version=version+1 and stocks =stocks -1 where version =version and goods=goods
			if result := tx.Model(&model.Inventory{}).Where("version = ? and goods =?", inv.Version, inv.Goods).Select("version", "stocks").Updates(&model.Inventory{
				Stocks:  inv.Stocks,
				Version: inv.Version + 1,
			}); result.RowsAffected == 0 {
				zap.S().Info("扣减库存失败...")
			} else {
				break
			}

			//if ok, err := mutex.Unlock(); !ok || err != nil {
			//	return errcode.ErrServer.WithDetails("释放redis分布式锁异常")
			//}
		}
		fmt.Println("over")
	}
	tx.Commit()
	return nil
}
func Init() *gorm.DB {
	dsn := "root:zxz123456@tcp(127.0.0.1:3306)/mall?charset=utf8mb4&parseTime=True&loc=Local"
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
	//_ = DB.AutoMigrate(&model.Inventory{})
	//for i := 421; i <= 600; i++ {
	//	var inv model.Inventory
	//	inv.Goods = int32(i)
	//	inv.Stocks = 100
	//	inv.Version = 0
	//	DB.Create(&inv)
	//}
	return DB
}
func main() {
	db := Init()
	req := form.SellInfo{
		GoodsInfo: []*form.GoodInfo{
			{421, 1},
		},
	}
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go Sell(&wg, db, req)
	}
	wg.Wait()
}
