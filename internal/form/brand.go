package form

type CreateBrand struct {
	Name string `json:"Name" form:"Name" binding:"required"`
	Logo string `json:"Logo" form:"logo" binding:"required"`
}

type UpdateBrand struct {
	Name string `json:"Name" form:"Name"`
	Logo string `json:"Logo" form:"logo"`
}

type ReqBrandsList struct {
	Page       int64 `json:"Page"`
	PagePerNum int64 `json:"PagePerNum"`
}
