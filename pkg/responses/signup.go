package responses

type User struct {
	UserId   string `json:"user_id"`
	Nickname string `json:"nickname"`
}

type SignUpResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}
