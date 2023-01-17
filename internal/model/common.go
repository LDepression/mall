package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"mall/internal/pkg/token"
	"time"
)

type PalLoad struct {
	PalLoad token.Payload
	Role    int
}
type BaseModel struct {
	ID        int32          `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"column:add_time"`
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_time"`
	IsDelete  bool           `gorm:"column:is_deleted"`
}
type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
