package query

type group struct {
	User          user
	Good          good
	Category      category
	Brand         brand
	CategoryBrand categoryBrand
	Inventory     inventory
	ShopCart      shopCart
	Order         order
}

var Group = new(group)
