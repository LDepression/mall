package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 ` gorm:"type:int;column:stocks;default:0;not null"`
	Version int32 ` gorm:"type:int;default:0;not null"` //涉及分布式锁的乐观锁
}
