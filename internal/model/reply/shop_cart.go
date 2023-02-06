package reply

type ShopCartInfoResponse struct {
	ID      int32 `json:"ID"`
	UserID  int32 `json:"UserID"`
	GoodID  int32 `json:"GoodID"`
	Num     int32 `json:"Num"`
	Checked bool  `json:"Checked"`
}

type CartItemListResponse struct {
	Total int64                   `json:"Total"`
	Data  []*ShopCartInfoResponse `json:"Data"`
}
