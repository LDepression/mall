package query

import (
	"errors"
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/model"
	"mall/internal/model/reply"
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
func (c *category) CreateCategory(ModelCategory model.Category) error {
	result := dao.Group.DB.Create(&ModelCategory)
	return result.Error
}
func (c *category) DeleteCategory(id int32) error {
	result := dao.Group.DB.Delete(&model.Category{}, id)
	return result.Error
}

func (c *category) UpdateCategory(id int, cMap map[string]interface{}) error {
	result := dao.Group.DB.Model(&model.Category{}).Where("id=?", id).Updates(cMap)
	return result.Error
}

func (c *category) SearchCategory(id int) (replyInfo reply.CategoryInfo, err error) {

	var basicInfos []reply.CategoryBasicInfo
	var SelfInfo reply.CategoryBasicInfo
	var category model.Category
	result := dao.Group.DB.Find(&category, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return replyInfo, err
	}
	replyInfo.Total = int32(result.RowsAffected)
	SelfInfo.ID = int32(id)
	SelfInfo.ParentCategoryID = category.ParentCategoryID
	SelfInfo.Name = category.Name
	SelfInfo.Level = category.Level

	var subCatgories []model.Category
	//将商品的子分类查询出来

	//将子分类查询出来

	dao.Group.DB.Where("parent_category_id=?", id).Find(&subCatgories)
	for _, subCategory := range subCatgories {
		basicInfos = append(basicInfos, reply.CategoryBasicInfo{
			ID:               subCategory.ID,
			Name:             subCategory.Name,
			ParentCategoryID: int32(id),
			Level:            subCategory.Level,
		})
	}
	replyInfo.SubCategories = basicInfos
	replyInfo.CategoryBasicInfo = SelfInfo

	return replyInfo, nil
}

func (c *category) GetAllCategoryList() ([]model.Category, error) {
	var categories []model.Category
	result := dao.Group.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categories)
	return categories, result.Error
}

func (c *category) CheckCategoryName(Name string) bool {

	result := dao.Group.DB.Where("name=?", Name).Find(&model.Category{})
	if result.RowsAffected == 0 { //说明没有这个分类
		return true
	} else {
		return false
	}
}
