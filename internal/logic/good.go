package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"mall/internal/dao"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/global"
	"mall/internal/model"
	"mall/internal/model/reply"
	"mall/internal/pkg/app/errcode"
)

func Model2Response(good model.Good) reply.RepGoodInfo {
	return reply.RepGoodInfo{
		ID: good.ID,
		Category: reply.ReqCategoryInfo{
			ID:   good.Category.ID,
			Name: good.Category.Name,
		},
		Brand: reply.ReqBrandInfo{
			ID:   good.Brand.ID,
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

func Req2Model(req form.CreateGoodReq) {

}

type good struct {
}

func (good) GetGoodsList(req form.GoodsFilterRequest) (reply.RepGoodsInfo, errcode.Err) {
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

//func (good) CreateGood(req form.CreateGoodReq) (reply.RepGoodsInfo, errcode.Err) {
//	var replyGoodsInfo reply.RepGoodsInfo
//	Qgood := query.NewGood()
//	var good model.Good
//	good.Name = req.Name
//	goodInfo, err := Qgood.GetGoodByName(good)
//	if err == nil {
//		//此时说明商品已经存在了
//		return replyGoodsInfo, errcode.ErrGoodExsit
//	}
//
//	Qgood.CreateGood()
//}
