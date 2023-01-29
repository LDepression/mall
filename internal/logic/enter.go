package logic

type group struct {
	User          user
	Email         pemail
	Upload        upload
	Good          good
	Category      category
	Brand         brand
	CategoryBrand categoryBrand
}

var Group = new(group)
