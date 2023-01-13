package logic

type group struct {
	User   user
	Email  pemail
	Upload upload
	Good   good
}

var Group = new(group)
