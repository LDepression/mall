package form

type CreateCategoryBrand struct {
	CategoryID int32 `json:"CategoryID" form:"CategoryID" binding:"required"`
	BrandID    int32 `json:"BrandID" form:"BrandID" binding:"required"`
}

type CategoryBrandList struct {
	Page       int `json:"Page" `
	PagePerNum int `json:"PagePerNum"`
}

type UpdateCategoryBrand struct {
	CategoryID int32 `json:"CategoryID" form:"CategoryID"`
	BrandID    int32 `json:"BrandID" form:"BrandID"`
}
