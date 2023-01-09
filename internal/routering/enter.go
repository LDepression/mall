package routering

type group struct {
	User  user
	Base  abase
	Email email
}

var Group = new(group)
