/**
 * @Author: lenovo
 * @Description:
 * @File:  limit_test
 * @Version: 1.0.0
 * @Date: 2023/03/01 20:38
 */

package query

import (
	"context"
	"mall/internal/dao"
	"testing"

	"github.com/go-redis/redis/v8"
)

var luaScript = ""

func TestLimit(t *testing.T) {
	redis.NewScript(`
	
`)
	dao.Group.Redis.rdb.Eval(context.Background())
}
