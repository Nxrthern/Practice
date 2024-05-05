package requests

type SignUpRequest struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}
