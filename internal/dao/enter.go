package dao

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type group struct {
	DB    *gorm.DB
	Redis *redis.Client
}

var Group = new(group)
