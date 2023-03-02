/**
 * @Author: lenovo
 * @Description:
 * @File:  common
 * @Version: 1.0.0
 * @Date: 2023/02/26 14:38
 */

package routering

import (
	"mall/internal/config"
	limit "mall/internal/pkg/limiter/api"

	"golang.org/x/time/rate"
)

func GetLimiters(buckets []config.Bucket) limit.RateLimiter {
	limiters := make([]limit.RateLimiter, len(buckets))
	for i := range limiters {
		//这里是因为NewLimiter中创建的对象实现了Wait和Limit方法
		limiters[i] = rate.NewLimiter(limit.Per(buckets[i].Count, buckets[i].Duration), buckets[i].Burst)
		//第一个参数表示Limit 也就是duration的时间产生count个token
		//第二个参数表示Token桶容量的大小
	}
	//将这些限流器进行混合,然后再将颗粒小的放到前面
	return limit.MultiLimiter(limiters...)
}
