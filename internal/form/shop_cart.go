package form

type ShopCartItemForm struct {
	GoodsId int32 `json:"goods" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}

type ShopCartItemUpdateForm struct {
	Nums    int32 `json:"nums" binding:"required,min=1"`
	Checked *bool `json:"checked"`
}

type CartItemRequest struct {
	ID      int32 `json:"ID"`
	UserID  int32 `json:"UserID"`
	GoodID  int32 `json:"GoodID"`
	Nums    int32 `json:"nums"`
	Checked bool  `json:"Checked"`
}
