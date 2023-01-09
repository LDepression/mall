package form

type SendEmailCode struct {
	Email string `json:"email" binding:"email"`
}
