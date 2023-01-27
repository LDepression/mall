package routering

type group struct {
	User     user
	Base     abase
	Email    email
	Upload   upload
	Good     good
	Category category
	Brand    brand
}

var Group = new(group)
