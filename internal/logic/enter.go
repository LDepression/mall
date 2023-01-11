package logic

type group struct {
	User   user
	Email  pemail
	Upload upload
}

var Group = new(group)
