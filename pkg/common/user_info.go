package common

type UserInfo struct {
	UserId   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Comment  string `json:"comment"`
	Password string `json:"password"`
}
