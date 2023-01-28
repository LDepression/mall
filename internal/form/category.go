package form

type CreateCategory struct {
	Name             string `json:"name" form:"name" binding:"required"`
	Level            int    `json:"level" form:"level" binding:"required"`
	ParentCategoryID int32  `json:"parent_category_id" form:"parent_category_id" binding:"required"`
	IsTab            bool   `json:"isTab" form:"isTab" binding:"required"`
}

type UpdateCategory struct {
	Name  string `form:"name" json:"name" binding:"required,min=3,max=20"`
	IsTab bool   `form:"is_tab" json:"is_tab"`
}
