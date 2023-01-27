package reply

type ReqCategoryInfo struct {
	ID   int32  `json:"ID"`
	Name string `json:"name"`
}

type CategoryInfo struct {
	CategoryBasicInfo CategoryBasicInfo   `json:"CategoryBasicInfo"`
	SubCategories     []CategoryBasicInfo `json:"SubCategories"`
	Total             int32               `json:"Total"`
}

type CategoryBasicInfo struct {
	ID               int32  `json:"ID"`
	Name             string `json:"name"`
	ParentCategoryID int32  `json:"parentCategoryID"`
	Level            int32  `json:"level"`
}

type AllCategoryData struct {
	JsonData string `json:"JsonData"`
}
