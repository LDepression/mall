package query

import (
	"mall/internal/dao"
	"mall/internal/model"

	"gorm.io/gorm"
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
func (good) CheckGoodByName(good model.Good) (bool, error) {
	var good1 model.Good
	result := dao.Group.DB.Where("name=?", good.Name).First(&good1)
	if result.RowsAffected == 0 {
		//说明没有这个
		return false, result.Error
	} else {
		return true, result.Error
	}
}

func (g good) CheckGoodByID(id int32) bool {
	if result := dao.Group.DB.Where("id=?", id).Find(&model.Good{}); result.RowsAffected == 0 {
		return false
	}
	return true
}

func (good) CreateGood(goodInfo model.Good) (model.Good, error) {
	result := dao.Group.DB.Model(&model.Good{}).Save(&goodInfo)
	return goodInfo, result.Error
}

func (g good) DeleteGood(id int32) error {
	result := dao.Group.DB.Delete(&model.Good{BaseModel: model.BaseModel{ID: id}})
	return result.Error
}

func (g good) GetGoodByID(id int32) (*model.Good, error) {
	var good model.Good
	if result := dao.Group.DB.Find(&good, id); result.RowsAffected == 0 {
		return nil, result.Error
	} else {
		return &good, nil
	}
}

func (g good) BatchGetGood(ids []int32) ([]model.Good, error) {
	var goods []model.Good
	if result := dao.Group.DB.Find(&goods, ids); result.RowsAffected == 0 {
		return goods, result.Error
	} else {
		return goods, nil
	}

}
