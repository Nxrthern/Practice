package responses

import (
	"practice/pkg/common"
)

type GetUserInfoConf struct {
	Message string          `json:"message"`
	User    common.UserInfo `json:"user"`
}

type GetUserInfo struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}
