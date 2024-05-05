package responses

type CommNick struct {
	Nickname string `json:"nickname"`
	Comment  string `json:"comment"`
}

type PatchUserInfo struct {
	Message string     `json:"message"`
	Recipe  []CommNick `json:"recipe"`
}
