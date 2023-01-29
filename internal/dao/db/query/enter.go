package query

type group struct {
	User          user
	Good          good
	Category      category
	Brand         brand
	CategoryBrand categoryBrand
}

var Group = new(group)
