package reply

type OrderInfoResponse struct {
	Id          int32
	UserId      int32
	Name        string
	Address     string
	Post        string
	OrderSn     string
	OrderStatus string
	AddTime     string
	OrderType   string
	Total       float32
}

type OrderListResponse struct {
	Total int64                `json:"total"`
	Data  []*OrderInfoResponse `json:"data"`
}

type OrderItemResponse struct {
	Id        int32   `json:"id"`
	OrderID   int32   `json:"order_id"`
	UserId    int32   `json:"user_id"`
	GoodID    int32   `json:"good_id"`
	GoodName  string  `json:"good_name"`
	GoodImage string  `json:"good_image"`
	Price     float32 `json:"price"`
	Num       int32   `json:"num"`
}

type OrderDetailsResponse struct {
	Total     int32 `json:"total"`
	GoodData  []*OrderItemResponse
	OrderData *OrderInfoResponse
}
