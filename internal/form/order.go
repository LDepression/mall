package form

type OrderFilterRequest struct {
	UserID     int64 `json:"user_id"`
	Page       int64 `json:"page" binding:"required"`
	PagePerNum int64 `json:"pagePerNum" binding:"required"`
}

type OrderRequest struct {
	ID      int32  `json:"id"`
	UserID  int32  `json:"user_id"`
	Address string `json:"address" binding:"required"`
	Post    string `json:"post"`
	Name    string `json:"name"`
	Mobile  string `json:"mobile" binding:"required"`
}

type OrderStatus struct {
	UserID  int64  `json:"user_id"`
	OrderSn string `json:"order_sn" binding:"required"`
	Status  string `json:"status" binding:"required"`
}
