/**
 * @Author: lenovo
 * @Description:
 * @File:  limit
 * @Version: 1.0.0
 * @Date: 2023/02/26 15:41
 */

package query

import (
	"context"
	"errors"

	"mall/internal/pkg/utils"
)

var ErrParse = errors.New("解析错误")

const limit = "return redis.call('CL.THROTTLE', KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5])"

type BucketRequest struct {
	Key     string
	Cap     int64 // 令牌桶容量-1
	GenNum  int64 // 令牌产生数,(与下一个参数一起)表示在指定的时间内令牌桶允许访问的次数
	GenTime int64 // 令牌产生时间,指定的时间
	Cost    int64 // 本次取走的令牌数
}

type BucketReply struct {
	Success  bool  // true 成功
	Capital  int64 // 令牌桶容量
	Buckets  int64 // 剩余令牌数(可用令牌)
	WaitTime int64 // -1/等待时间(若请求被拒绝,这个值表示多久后令牌桶会重新添加令牌)
	FullTime int64 // 预计多少秒会满
}

func (q *Queries) GetBucket(ctx context.Context, config BucketRequest) (*BucketReply, error) {
	key := config.Key
	capital := utils.IDToSting(config.Cap)
	genNum := utils.IDToSting(config.GenNum)
	genTime := utils.IDToSting(config.GenTime)
	cost := utils.IDToSting(config.Cost)
	ret, err := q.rdb.Eval(ctx, limit, []string{key, capital, genNum, genTime, cost}).Result()
	if err != nil {
		return nil, err
	}
	res, ok := ret.([]interface{})
	if !ok || len(res) != 5 {
		return nil, ErrParse
	}
	return &BucketReply{
		Success:  res[0].(int64) == 0,
		Capital:  res[1].(int64),
		Buckets:  res[2].(int64),
		WaitTime: res[3].(int64),
		FullTime: res[4].(int64),
	}, nil
}
