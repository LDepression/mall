package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"gorm.io/gorm/utils"
	"mall/internal/dao"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/global"
	"mall/internal/model"
	"mall/internal/model/reply"
	"mall/internal/pkg/app/errcode"
	"sync"
)

func Model2Response(good model.Good) reply.RepGoodInfo {
	return reply.RepGoodInfo{
		ID: good.ID,
		Category: reply.ReqCategoryInfo{
			ID:   good.CategoryID,
			Name: good.Category.Name,
		},
		Brand: reply.ReqBrandInfo{
			ID:   good.BrandID,
			Name: good.Brand.Name,
			Logo: good.Brand.Logo,
		},
		OnSale:          good.OnSale,
		ShipFree:        good.ShipFree,
		IsNew:           good.IsNew,
		IsHot:           good.IsHot,
		Name:            good.Name,
		GoodsSn:         good.GoodsSn,
		ClickNum:        good.ClickNum,
		SoldNum:         good.SoldNum,
		FavNum:          good.FavNum,
		MarketPrice:     good.MarketPrice,
		ShopPrice:       good.ShopPrice,
		GoodsBrief:      good.GoodsBrief,
		Images:          good.Images,
		DescImages:      good.DescImages,
		GoodsFrontImage: good.GoodsFrontImage,
	}
}

func Req2Model(req form.CreateGoodReq, category model.Category, brand model.Brand) model.Good {
	return model.Good{
		Category:        category,
		CategoryID:      req.CategoryId,
		BrandID:         req.Brand,
		Brand:           brand,
		ShipFree:        *req.ShipFree,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.FrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
	}
}

type good struct {
	Lock *sync.Mutex
}

func (g *good) GetGoodsList(req form.GoodsFilterRequest) (reply.RepGoodsInfo, errcode.Err) {
	//我们使用es进行查询
	/*
		mysql作为存储,es进行查询
	*/
	var ReplyGoodInfo reply.RepGoodsInfo
	q := elastic.NewBoolQuery()
	if req.KeyWords != "" {
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	//filter不参与算分
	if req.IsHot {
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}
	if req.IsTab {
		q = q.Filter(elastic.NewTermQuery("is_tab", req.IsTab))
	}
	if req.PriceMax != 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.PriceMin != 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMax))
	}

	//还可以根据category进行查询
	var subCategory string
	var categoryIDs []int
	//topCategory就是商品的categoryID
	if req.TopCategory != 0 {
		qCategory := query.NewCategory()
		categoryInfo, err := qCategory.GetCategoryByID(req.TopCategory)
		if err != nil {
			zap.S().Info("没有该分类")
			return ReplyGoodInfo, errcode.ErrServer.WithDetails(err.Error())
		}
		if categoryInfo.Level == 1 {
			subCategory = fmt.Sprintf("select id from categories where parent_category_id in(select id from categories where parent_category_id=%d)", req.TopCategory)
		} else if categoryInfo.Level == 2 {
			subCategory = fmt.Sprintf("select id from categories where parent_category_id=%d", req.TopCategory)
		} else if categoryInfo.Level == 3 {
			subCategory = fmt.Sprintf("select id from categories where id=%d", req.TopCategory)
		}
		//subCategory查询出来的就是
		type Result struct {
			ID int `json:"id"`
		}
		var results []Result
		//将id扫描到result中去
		dao.Group.DB.Model(&model.Category{}).Raw(subCategory).Scan(&results)
		for _, re := range results {
			categoryIDs = append(categoryIDs, re.ID)
		}
		q.Filter(elastic.NewTermsQuery("category_id", categoryIDs))
	}
	//关于分页的操作
	if req.Pages == 0 {
		req.Pages = 0
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	//如果个数不够的话,那么如果我们强行增加pages和pagepernum的话,就查询不出来数据
	result, err := global.EsClient.Search().Index(model.EsGoods{}.IndexName()).Query(q).Do(context.Background())
	if err != nil {
		return ReplyGoodInfo, errcode.ErrServer.WithDetails(err.Error())
	}
	ReplyGoodInfo.Total = int32(result.Hits.TotalHits.Value)
	if result.Hits.TotalHits.Value == 0 {
		return ReplyGoodInfo, errcode.ErrNotFound
	}
	var goodIDs []int32
	for _, v := range result.Hits.Hits {
		good := model.EsGoods{}
		_ = json.Unmarshal(v.Source, &good)
		goodIDs = append(goodIDs, good.ID)
	}
	//就是通过goodIDs批量得去查询商品,
	var goods []model.Good
	/*
		值得注意的是,要是什么goodIDs都没有的话,就会全部查询出来
	*/
	dao.Group.DB.Model(&model.Good{}).Preload("Category").Preload("Brand").Scopes(query.Paginate(int64(req.Pages), int64(req.PagePerNums))).Find(&goods, goodIDs)
	for _, v := range goods {
		rlp := Model2Response(v)
		ReplyGoodInfo.Data = append(ReplyGoodInfo.Data, rlp)
	}
	return ReplyGoodInfo, nil
}

func (g *good) CreateGood(ctx *gin.Context, req form.CreateGoodReq) (*reply.RepGoodInfo, errcode.Err) {
	var category model.Category
	if result := dao.Group.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, errcode.ErrServer.WithDetails("分类不存在")
	}
	var brand model.Brand
	if result := dao.Group.DB.First(&brand, req.Brand); result.RowsAffected == 0 {
		return nil, errcode.ErrServer.WithDetails("商品的品牌不存在")
	}
	Qgood := query.NewGood()
	good := Req2Model(req, category, brand)
	exist, _ := Qgood.CheckGoodByName(good)
	if exist {
		return nil, errcode.ErrGoodExsit
	}
	//创建商品
	tx := Qgood.Begin()
	//_, err := Qgood.CreateGood(good)
	result := tx.Save(&good)
	if result.Error != nil {
		tx.Rollback()
		return nil, errcode.ErrServer
	}
	tx.Commit()
	res := Model2Response(good)
	//这里设置进redis缓存中去
	err1 := dao.Group.Redis.SetGood(ctx, res)
	if err1 != nil {
		return nil, errcode.ErrServer.WithDetails(err1.Error())
	}
	return &res, nil
}

func (g *good) DeleteGood(ctx *gin.Context, id int32) errcode.Err {
	Qgood := query.NewGood()
	if exsit := Qgood.CheckGoodByID(id); !exsit {
		return errcode.ErrNotFound
	}
	tx := Qgood.Begin()
	result := tx.Delete(&model.Good{BaseModel: model.BaseModel{ID: id}})
	if result.Error != nil {
		tx.Rollback()
		return errcode.ErrServer.WithDetails(result.Error.Error())
	}
	tx.Commit()
	g.Lock.Lock()
	defer g.Lock.Unlock()
	err1 := dao.Group.Redis.Del(ctx, utils.ToString(id))
	if err1 != nil {
		return errcode.ErrServer.WithDetails(err1.Error())
	}
	return nil
}

func (g *good) UpdateGood(ctx *gin.Context, updateGood form.UpdateGood) (*reply.RepGoodInfo, errcode.Err) {
	//保证一致性
	g.Lock.Lock()
	defer g.Lock.Unlock()
	var category model.Category
	if result := dao.Group.DB.First(&category, updateGood.CategoryId); result.RowsAffected == 0 {
		return nil, errcode.ErrServer.WithDetails("分类不存在")
	}
	var brand model.Brand
	if result := dao.Group.DB.First(&brand, updateGood.Brand); result.RowsAffected == 0 {
		return nil, errcode.ErrServer.WithDetails("商品的品牌不存在")
	}

	var replyGood reply.RepGoodInfo
	var good model.Good
	Qgood := query.NewGood()
	if exsit := Qgood.CheckGoodByID(updateGood.ID); !exsit {
		return nil, errcode.ErrNotFound
	}
	good.ID = updateGood.ID
	good.Name = updateGood.Name
	good.GoodsSn = updateGood.GoodsSn
	good.GoodsBrief = updateGood.GoodsBrief
	good.GoodsFrontImage = updateGood.FrontImage
	good.BrandID = updateGood.Brand
	good.CategoryID = updateGood.CategoryId
	good.MarketPrice = updateGood.MarketPrice
	good.ShopPrice = updateGood.ShopPrice
	good.IsHot = updateGood.IsHot
	good.IsNew = updateGood.IsNew
	good.Category = category
	good.Brand = brand
	replyGood = Model2Response(good)
	tx := dao.Group.DB.Begin()
	result := tx.Updates(&good)
	if result.Error != nil {
		tx.Rollback()
		return nil, errcode.ErrServer.WithDetails(result.Error.Error())
	}
	tx.Commit()
	if err := dao.Group.Redis.SetGood(ctx, replyGood); err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	return &replyGood, nil
}

func (g *good) GetGoodByID(ctx *gin.Context, id int32) (*reply.RepGoodInfo, errcode.Err) {
	replyInfo, err := dao.Group.Redis.GetGood(ctx, int64(id))
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			zap.S().Error("redis 内部错误....")
		}
	} else {
		return &replyInfo, nil
	}
	Qgood := query.NewGood()
	goodInfo, err := Qgood.GetGoodByID(id)
	if err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	reply1 := Model2Response(*goodInfo)
	return &reply1, nil
}

func (g *good) BatchGetGoods(ids []int) ([]reply.RepGoodInfo, errcode.Err) {
	var replyGoodInfo []reply.RepGoodInfo
	Qgood := query.NewGood()
	goodsInfo, err := Qgood.BatchGetGood(ids)
	if err != nil {
		return replyGoodInfo, errcode.ErrServer.WithDetails(err.Error())
	}
	for _, goodInfo := range goodsInfo {
		t := Model2Response(goodInfo)
		replyGoodInfo = append(replyGoodInfo, t)
	}
	return replyGoodInfo, nil
}
