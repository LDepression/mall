package query

type group struct {
	User     user
	Good     good
	Category category
	Brand    brand
}

var Group = new(group)
