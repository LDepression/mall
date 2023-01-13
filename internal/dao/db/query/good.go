package query

import (
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/model"
)

type good struct {
	*gorm.DB
}

func NewGood() *good {
	return &good{
		dao.Group.DB,
	}
}

func (good) GetGoodByName(good2 model.Good) (model.Good, error) {
	result := dao.Group.DB.Where("name=?", good2.Name).First(&good2)
	return good2, result.Error
}

//func (good) CreateGood(goodInfo form.CreateGoodReq) (reply.RepGoodInfo, error) {
//
//}
