package form

import "mall/internal/model"

type GoodsFilterRequest struct {
	PriceMin    int32  `json:"priceMin"`
	PriceMax    int32  `json:"priceMax"`
	IsHot       bool   `json:"isHot"`
	IsNew       bool   `json:"isNew"`
	IsTab       bool   `json:"isTab"`
	TopCategory int32  `json:"topCategory"` //主键id
	Pages       int32  `json:"pages"`
	PagePerNums int32  `json:"pagePerNums"`
	KeyWords    string `json:"keyWords"`
	Brand       int32  `json:"brand"`
}

type CreateGoodReq struct {
	Name        string         `form:"name" json:"name" binding:"required,min=2,max=100"`
	GoodsSn     string         `form:"goods_sn" json:"goods_sn" binding:"required,min=2,lt=20"`
	Stocks      int32          `form:"stocks" json:"stocks" binding:"required,min=1"`
	CategoryId  int32          `form:"category" json:"category" binding:"required"`
	MarketPrice float32        `form:"market_price" json:"market_price" binding:"required,min=0"`
	ShopPrice   float32        `form:"shop_price" json:"shop_price" binding:"required,min=0"`
	GoodsBrief  string         `form:"goods_brief" json:"goods_brief" binding:"required,min=3"`
	Images      model.GormList `form:"images" json:"images" binding:"required,min=1"`
	DescImages  model.GormList `form:"desc_images" json:"desc_images" binding:"required,min=1"`
	ShipFree    *bool          `form:"ship_free" json:"ship_free" binding:"required"`
	FrontImage  string         `form:"front_image" json:"front_image" binding:"required,url"`
	Brand       int32          `form:"brand" json:"brand" binding:"required"`
	IsNew       bool           `form:"is_new" json:"is_new" binding:"required"`
	IsHot       bool           `json:"is_hot" json:"is_hot" binding:"required"`
	OnSale      bool           `json:"on_sale" json:"on_sale" binding:"required"`
}

type DeleteGoodReq struct {
	GoodID int32 `form:"goodID" json:"goodID" binding:"required"`
}

type UpdateGood struct {
	ID          int32          `form:"id" json:"id"`
	Name        string         `form:"name" json:"name" `
	GoodsSn     string         `form:"goods_sn" json:"goods_sn" `
	Stocks      int32          `form:"stocks" json:"stocks" `
	CategoryId  int32          `form:"category" json:"category" binding:"required"`
	MarketPrice float32        `form:"market_price" json:"market_price" `
	ShopPrice   float32        `form:"shop_price" json:"shop_price" `
	GoodsBrief  string         `form:"goods_brief" json:"goods_brief" `
	Images      model.GormList `form:"images" json:"images" `
	DescImages  model.GormList `form:"desc_images" json:"desc_images"`
	ShipFree    *bool          `form:"ship_free" json:"ship_free" `
	FrontImage  string         `form:"front_image" json:"front_image" `
	Brand       int32          `form:"brand" json:"brand" binding:"required"`
	IsNew       bool           `form:"is_new" json:"is_new" `
	IsHot       bool           `json:"is_hot" json:"is_hot" `
	OnSale      bool           `json:"on_sale" json:"on_sale" `
}

type GetGoodID struct {
	ID int32 `json:"id" binding:"required" form:"id"`
}
