package model

import (
	"database/sql/driver"
	"encoding/json"
)

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 ` gorm:"type:int;column:stocks;default:0;not null"`
	Version int32 ` gorm:"type:int;default:0;not null"` //涉及分布式锁的乐观锁
}
type GoodInfo struct {
	GoodID int32 `gorm:"type:int"`
	Num    int32 `gorm:"type:int"`
}
type RebackGoodList []GoodInfo
type RebackDetails struct {
	GoodsInfo RebackGoodList
	OrderSn   string `gorm:"type:varchar(200);index:idx_order_sn,unique"`
	Status    int32  `gorm:"type:int"` //1表示已支付,2表示已归还
}

func (RebackDetails) TableName() string {
	return "rebackdetails"
}
func (g RebackGoodList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *RebackGoodList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
