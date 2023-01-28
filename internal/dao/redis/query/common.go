package query

import (
	"context"
	"encoding/json"
	"mall/internal/global"
	"mall/internal/pkg/singleflight"
	"time"
)

// Set 设置数据
func (q *Queries) Set(ctx context.Context, key string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return q.rdb.Set(ctx, key, data, global.Setting.Redis.CacheTime).Err()
}

//设置具有特定过期时间的数据

func (q *Queries) SetTimeOut(ctx context.Context, key string, val interface{}, duration time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return q.rdb.Set(ctx, key, data, duration).Err()
}

// Get 获取数据绑定到val上 val需要指针形式
func (q *Queries) Get(ctx context.Context, key string, val interface{}) error {
	result, err := singleflight.Group.Do(key, func() (interface{}, error) {
		result, err := q.rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(result.(string)), val); err != nil {
		return err
	}
	return nil
}

func (q *Queries) Del(ctx context.Context, key string) error {
	_, err := singleflight.Group.Do(key, func() (interface{}, error) {
		return nil, q.rdb.Expire(ctx, key, 0).Err()
	})
	return err
}
