package logic

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/model"
	"mall/internal/model/reply"
	"mall/internal/pkg/app/errcode"
	"sync"
)

type categoryBrand struct {
	Lock *sync.Mutex
}

func (b *categoryBrand) CreateCategoryBrand(form form.CreateCategoryBrand) (*reply.CategoryBrandResponse, errcode.Err) {
	var rsp reply.CategoryBrandResponse
	qCategoryBrand := query.NewCategoryBrand()
	rspInfo, err := qCategoryBrand.CreateCategoryBrand(form.CategoryID, form.BrandID)
	if err != nil {
		zap.S().Info("qCategoryBrand.CreateCategoryBrand failed,err:", err)
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	rsp.ID = int(rspInfo.ID)
	return &rsp, nil
}
func (b *categoryBrand) DeleteCategoryBrand(id int) errcode.Err {
	qCategoryBrand := query.NewCategoryBrand()
	if err := qCategoryBrand.DeleteCategoryBrand(id); err != nil {
		zap.S().Info("qCategoryBrand.DeleteCategoryBrand(id) failed,err:", err)
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func (b *categoryBrand) UpdateCategoryBrand(id int, form form.UpdateCategoryBrand) errcode.Err {
	qCategoryBrand := query.NewCategoryBrand()
	err := qCategoryBrand.UpdateCategoryBrand(id, model.GoodsCategoryBrand{
		CategoryID: form.CategoryID,
		BrandID:    form.BrandID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func (b *categoryBrand) CategoryBrandList(page, pagePerNum int) (*reply.CategoryBrand, errcode.Err) {
	rsp := new(reply.CategoryBrand)
	var CategoryBrand []*model.GoodsCategoryBrand
	//先来查询总数
	result := dao.Group.DB.Find(&CategoryBrand)
	rsp.Total = result.RowsAffected
	result = dao.Group.DB.Scopes(query.Paginate(int64(page), int64(pagePerNum))).Preload("Category").Preload("Brand").Find(&CategoryBrand)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errcode.ErrNotFound
	}
	for _, categoryBrandInfo := range CategoryBrand {
		rsp.Data = append(rsp.Data, &reply.CategoryBrandResponse{
			ID: int(categoryBrandInfo.ID),
			Brand: reply.BrandResponse{
				ID:   int(categoryBrandInfo.Brand.ID),
				Name: categoryBrandInfo.Brand.Name,
				Logo: categoryBrandInfo.Brand.Logo,
			},
			Category: reply.CategoryResponse{
				ID:    int(categoryBrandInfo.Category.ID),
				Name:  categoryBrandInfo.Category.Name,
				Level: int(categoryBrandInfo.Category.Level),
			},
		})
	}
	return rsp, nil
}
