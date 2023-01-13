package reply

type RepGoodsInfo struct {
	Data  []RepGoodInfo
	Total int32 `json:"total"`
}

type RepGoodInfo struct {
	ID int32 `json:"id"`

	Category ReqCategoryInfo `json:"category"`
	Brand    ReqBrandInfo    `json:"brand"`

	OnSale   bool `json:"onSale"`
	ShipFree bool `json:"shipFree"`
	IsNew    bool `json:"isNew"`
	IsHot    bool `json:"isHot"`

	Name            string   `json:"name"`
	GoodsSn         string   `json:"goodsSn"`
	ClickNum        int32    `json:"clickNum"`
	SoldNum         int32    `json:"soldNum"`
	FavNum          int32    `json:"favNum"`
	MarketPrice     float32  `json:"marketPrice"`
	ShopPrice       float32  `json:"shopPrice"`
	GoodsBrief      string   `json:"goodsBrief"`
	Images          []string `json:"images"`
	DescImages      []string `json:"descImages"`
	GoodsFrontImage string   `json:"goodsFrontImage"`
}
