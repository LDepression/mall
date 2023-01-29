package query

import (
	"mall/internal/dao"
	"mall/internal/model"
)

type categoryBrand struct {
}

func NewCategoryBrand() *categoryBrand {
	return &categoryBrand{}
}

func (b *categoryBrand) CreateCategoryBrand(categoryID, brandID int32) (*model.GoodsCategoryBrand, error) {
	var categoryBrand model.GoodsCategoryBrand
	categoryBrand.CategoryID = categoryID
	categoryBrand.BrandID = brandID
	result := dao.Group.DB.Create(&categoryBrand)
	return &categoryBrand, result.Error
}
func (b *categoryBrand) DeleteCategoryBrand(id int) error {
	result := dao.Group.DB.Delete(&model.GoodsCategoryBrand{}, id)
	return result.Error
}
func (b *categoryBrand) UpdateCategoryBrand(id int, goodsCategoryBrand model.GoodsCategoryBrand) error {
	result := dao.Group.DB.Where("id=?", id).Updates(&goodsCategoryBrand)
	return result.Error
}
