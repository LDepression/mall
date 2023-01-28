package logic

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/model"
	"mall/internal/model/reply"
	"mall/internal/pkg/app/errcode"
)

type category struct {
}

func (*category) CreateCategory(categoryForm form.CreateCategory) errcode.Err {
	qCategory := query.NewCategory()
	if ok := qCategory.CheckCategoryName(categoryForm.Name); !ok {
		//已经有了
		return errcode.ErrCategoryNameExist
	}

	//先去判断一下商品是否存在
	modelCategory := model.Category{
		Name:  categoryForm.Name,
		Level: int32(categoryForm.Level),
		IsTab: false,
	}
	if categoryForm.Level != 1 {
		modelCategory.ParentCategoryID = categoryForm.ParentCategoryID
	}
	if err := qCategory.CreateCategory(modelCategory); err != nil {
		zap.S().Info("qCategory.CreateCategory(cMap) failed", err)
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}
func (*category) DeleteCategory(categoryID int32) errcode.Err {
	qCategory := query.NewCategory()
	_, err := qCategory.GetCategoryByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return errcode.ErrServer.WithDetails(err.Error())
	}
	if err := qCategory.DeleteCategory(categoryID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func (*category) UpdateCategory(categoryForm form.UpdateCategory, categoryID int) errcode.Err {
	cMap := make(map[string]interface{})
	cMap["name"] = categoryForm.Name
	cMap["is_tab"] = categoryForm.IsTab
	qCategory := query.NewCategory()
	if err := qCategory.UpdateCategory(categoryID, cMap); err != nil {
		zap.S().Info("qCategory.UpdateCategory(cMap) failed", err)
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func (*category) SearchCategory(categoryID int) (reply.CategoryInfo, errcode.Err) {
	qCategory := query.NewCategory()
	categoryReply, err := qCategory.SearchCategory(categoryID)
	if err != nil {
		zap.S().Info("qCategory.UpdateCategory(cMap) failed", err)
		return categoryReply, errcode.ErrServer.WithDetails(err.Error())
	}
	//否则查询到了数据,此时返回数据就好了
	return categoryReply, nil
}
func (*category) GetAllCategoryList() (*reply.AllCategoryData, errcode.Err) {
	qCategory := query.NewCategory()
	categories, err := qCategory.GetAllCategoryList()
	if err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	b, _ := json.Marshal(&categories)
	return &reply.AllCategoryData{JsonData: string(b)}, nil
}
