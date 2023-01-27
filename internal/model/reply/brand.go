package reply

type ReqBrandInfo struct {
	ID   int32  `json:"ID"`
	Name string `json:"Name"`
	Logo string `json:"Logo"`
}
type RepBrandsList struct {
	BrandsInfo []ReqBrandInfo `json:"BrandsInfo"`
	Total      int32          `json:"Total"`
}
