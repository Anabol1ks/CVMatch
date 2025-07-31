package response

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserRegisterResponse struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}
