package query

import (
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/model"
)

type category struct {
	*gorm.DB
}

func NewCategory() *category {
	return &category{
		dao.Group.DB,
	}
}
func (c *category) GetCategoryByID(categoryID int32) (model.Category, error) {
	var categoryInfo model.Category
	result := c.First(&categoryInfo, categoryID)
	return categoryInfo, result.Error
}
