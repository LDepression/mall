package routering

type group struct {
	User   user
	Base   abase
	Email  email
	Upload upload
	Good   good
}

var Group = new(group)
