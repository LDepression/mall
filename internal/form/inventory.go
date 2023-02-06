package form

type GoodInfo struct {
	GoodID int32 `json:"GoodID"`
	Num    int32 `json:"Num"`
}

type SellInfo struct {
	GoodsInfo []*GoodInfo `json:"GoodsInfo"`
	OrderSn   string      `json:"OrderSn"`
}
