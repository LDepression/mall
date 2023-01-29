package reply

type CategoryBrand struct {
	Total int64                    `json:"Total"`
	Data  []*CategoryBrandResponse `json:"Data"`
}

type CategoryBrandResponse struct {
	ID       int              `json:"ID"`
	Brand    BrandResponse    `json:"Brand"`
	Category CategoryResponse `json:"Category"`
}

type BrandResponse struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
	Logo string `json:"Logo"`
}

type CategoryResponse struct {
	ID    int    `json:"ID"`
	Name  string `json:"Name"`
	Level int    `json:"Level"`
}
