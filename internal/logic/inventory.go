package logic

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"mall/internal/dao"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/model"
	"mall/internal/pkg/app/errcode"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type inventory struct {
}

//SetInv 设置库存
func (i *inventory) SetInv(req form.GoodInfo) errcode.Err {
	qInv := query.NewInventory()
	if err := qInv.CheckGoodID(req.GoodID); err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errcode.ErrNotFound
		}
		return errcode.ErrServer
	}
	var inv model.Inventory
	inv.Goods = req.GoodID
	inv.Stocks = req.Num
	if err := qInv.SetInv(req); err != nil {
		zap.S().Info("qInv.SetInv(req) failed,", err.Error())
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

//InvDetails 库存信息
func (i *inventory) InvDetails(req form.GoodInfo) (*form.GoodInfo, errcode.Err) {
	qInv := query.NewInventory()
	if err := qInv.CheckGoodID(req.GoodID); err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrServer
	}
	InvInfo, err := qInv.GetInvDetails(req.GoodID)
	if err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	return &form.GoodInfo{
		GoodID: InvInfo.Goods,
		Num:    InvInfo.Stocks,
	}, nil
}

//Sell 欲扣减
func (i *inventory) Sell(req form.SellInfo) errcode.Err {
	tx := dao.Group.DB.Begin()
	var sellDetails model.RebackDetails
	for _, goodInfo := range req.GoodsInfo {
		sellDetails = model.RebackDetails{
			OrderSn: req.OrderSn,
			Status:  1,
		}
		sellDetails.GoodsInfo = append(sellDetails.GoodsInfo, model.GoodInfo{
			GoodID: goodInfo.GoodID,
			Num:    goodInfo.Num,
		})
		for {
			var inv model.Inventory
			//先去查询一下
			if result := tx.Where("goods=?", goodInfo.GoodID).Find(&inv); result.RowsAffected == 0 {
				//说明此时没哟这个消息的记录,那么此时扣减库存就失败了
				tx.Rollback()
				return errcode.ErrServer.WithDetails(result.Error.Error())
			}
			//这里先去判断一下,商品的库存是否够用
			//这里表示的是如果此时库存不够
			if inv.Stocks < goodInfo.Num {
				tx.Rollback()
				return errcode.ErrServer.WithDetails("商品的库存不足....")
			}
			inv.Stocks -= goodInfo.Num //扣减库存
			//update inventory set version=version+1 and stocks =stocks -1 where version =version and goods=goods
			if result := tx.Model(&model.Inventory{}).Where("version = ? and goods =?", inv.Version, inv.Goods).Select("version", "stocks").Updates(&model.Inventory{
				Stocks:  inv.Stocks,
				Version: inv.Version + 1,
			}); result.RowsAffected == 0 {
				zap.S().Info("扣减库存失败...")
			} else {
				break
			}
		}
	}
	if result := tx.Create(&sellDetails); result.RowsAffected == 0 {
		return errcode.ErrServer
	}
	tx.Commit()
	return nil
}

//////Reback 库存归还
//func (i *inventory) Reback() errcode.Err {
//
//}

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
	DB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})
	return DB
}
func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range msgs {
		//既然是归还库存，那么我应该具体的知道每件商品应该归还多少， 但是有一个问题是什么？重复归还的问题
		//所以说这个接口应该确保幂等性， 你不能因为消息的重复发送导致一个订单的库存归还多次， 没有扣减的库存你别归还
		//如果确保这些都没有问题， 新建一张表， 这张表记录了详细的订单扣减细节，以及归还细节
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("解析json失败： %v\n", msgs[i].Body)
			return consumer.ConsumeSuccess, nil
		}
		zap.S().Info(msgs[i].Body)
		//去将inv的库存加回去 将selldetail的status设置为2， 要在事务中进行
		db := Init()
		tx := db.Begin()
		var sellDetail model.RebackDetails
		if result := tx.Model(&model.RebackDetails{}).Where(&model.RebackDetails{OrderSn: orderInfo.OrderSn, Status: 1}).First(&sellDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		//如果查询到那么逐个归还库存
		for _, orderGood := range sellDetail.GoodsInfo {
			//update怎么用
			//先查询一下inventory表在， update语句的 update xx set stocks=stocks+2
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: orderGood.GoodID}).Update("stocks", gorm.Expr("stocks+?", orderGood.Num)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}

		if result := tx.Model(&model.RebackDetails{}).Where(&model.RebackDetails{OrderSn: orderInfo.OrderSn}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
