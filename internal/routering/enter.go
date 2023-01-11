package routering

type group struct {
	User   user
	Base   abase
	Email  email
	Upload upload
}

var Group = new(group)
