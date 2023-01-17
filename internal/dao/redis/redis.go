package redis

import (
	"github.com/go-redis/redis/v8"
	"mall/internal/dao/redis/query"
	"mall/internal/global"
)

func Init() *query.Queries {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.Setting.Redis.Addr,
		Password: global.Setting.Redis.Password, // 密码
		DB:       0,                             // 数据库
		PoolSize: global.Setting.Redis.PoolSize, // 连接池大小
	})
	return query.New(rdb)
}
