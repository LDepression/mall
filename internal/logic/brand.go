package logic

import (
	"go.uber.org/zap"
	"mall/internal/dao"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/model"
	"mall/internal/model/reply"
	"mall/internal/pkg/app/errcode"
	"sync"
)

type brand struct {
	Lock *sync.Mutex
}

func (b *brand) CreateBrand(form form.CreateBrand) errcode.Err {

	//先去判断一下品牌是否存在
	qBrand := query.NewBrand()
	esxit := qBrand.CheckBrandByName(form.Name)
	if esxit {
		return errcode.ErrServer.WithDetails("品牌已经存在了")
	}
	var brand model.Brand
	brand.Name = form.Name
	brand.Logo = form.Logo
	b.Lock.Lock()
	defer b.Lock.Unlock()
	if err := qBrand.CreateBrand(brand); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}
func (b *brand) DeleteBrand(id int) errcode.Err {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	qBrand := query.NewBrand()
	if err := qBrand.DeleteBrand(id); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func (b *brand) UpdateBrand(id int, form form.UpdateBrand) errcode.Err {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	qBrand := query.NewBrand()
	brand := model.Brand{
		Name: form.Name,
		Logo: form.Logo,
	}
	if err := qBrand.UpdateBrand(id, brand); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func (b *brand) GetBrandByID(id int) (*reply.ReqBrandInfo, errcode.Err) {
	//先去判断一下是否存在
	qBrand := query.NewBrand()
	brandInfo, err := qBrand.GetBrandByID(id)
	if err != nil {
		zap.S().Info("qBrand.GetBrandByID(id) failed,err:", err)
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	var replyBrand reply.ReqBrandInfo
	replyBrand.ID = brandInfo.ID
	replyBrand.Name = brandInfo.Name
	replyBrand.Logo = brandInfo.Logo
	return &replyBrand, nil
}

func (*brand) BrandList(reqBrandList form.ReqBrandsList) (*reply.RepBrandsList, errcode.Err) {

	rep := reply.RepBrandsList{}
	var brandsInfo []reply.ReqBrandInfo
	var brands []model.Brand

	//查询出总数
	var Brands []model.Brand
	result := dao.Group.DB.Find(&Brands)
	if result.RowsAffected == 0 {
		return nil, errcode.ErrNotFound
	}
	rep.Total = int32(result.RowsAffected)
	result = dao.Group.DB.Scopes(query.Paginate(reqBrandList.Page, reqBrandList.PagePerNum)).Find(&brands)
	if result.RowsAffected == 0 {
		return nil, errcode.ErrNotFound
	}

	for _, brandInfo := range brands {
		brandsInfo = append(brandsInfo, reply.ReqBrandInfo{
			ID:   brandInfo.ID,
			Name: brandInfo.Name,
			Logo: brandInfo.Logo,
		})
	}
	rep.BrandsInfo = brandsInfo
	return &rep, nil
}
