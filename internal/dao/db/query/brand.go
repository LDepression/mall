package query

import (
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/model"
)

type brand struct {
	*gorm.DB
}

func NewBrand() *brand {
	return &brand{
		dao.Group.DB,
	}
}
func (b *brand) CreateBrand(brand model.Brand) error {
	result := dao.Group.DB.Create(&brand)
	return result.Error
}
func (b *brand) DeleteBrand(id int) error {
	result := dao.Group.DB.Delete(&model.Brand{}, id)
	return result.Error
}

func (b *brand) UpdateBrand(id int, brand model.Brand) error {
	result := dao.Group.DB.Where("id=?", id).Updates(&brand)
	return result.Error
}

func (b *brand) GetBrandByID(id int) (model.Brand, error) {
	var brand model.Brand
	result := dao.Group.DB.Find(&brand, id)
	return brand, result.Error
}
func (b *brand) CheckBrandByName(name string) bool {
	result := dao.Group.DB.Where(&model.Brand{Name: name}).Find(&model.Brand{})
	if result.RowsAffected == 1 {
		return true
	}
	return false
}
