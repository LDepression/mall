package routering

type group struct {
	User     user
	Base     abase
	Email    email
	Upload   upload
	Good     good
	Category category
	Brand    brand
	ShopCart shopCart
	Order    order
}

var Group = new(group)
