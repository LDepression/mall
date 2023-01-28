package query

import (
	"context"
	"github.com/0RAJA/Rutils/pkg/utils"
	utils2 "gorm.io/gorm/utils"
	"mall/internal/model/reply"
)

var KeyGood string = "KeyGood"

func (q *Queries) SetGood(ctx context.Context, good reply.RepGoodInfo) error {
	err := q.Set(ctx, utils.LinkStr(KeyGood, utils2.ToString(good.ID)), good)
	return err
}
func (q *Queries) DelGood(ctx context.Context, goodID int64) error {
	err := q.Del(ctx, utils.LinkStr(KeyGood, utils2.ToString(goodID)))
	return err
}
func (q *Queries) GetGood(ctx context.Context, goodID int64) (reply.RepGoodInfo, error) {
	var replyGoodInfo reply.RepGoodInfo
	err := q.Get(ctx, utils.LinkStr(KeyGood, utils2.ToString(goodID)), &replyGoodInfo)
	return replyGoodInfo, err
}
