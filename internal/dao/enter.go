package dao

import (
	"gorm.io/gorm"
	"mall/internal/dao/redis/query"
)

type group struct {
	DB    *gorm.DB
	Redis *query.Queries
}

var Group = new(group)
