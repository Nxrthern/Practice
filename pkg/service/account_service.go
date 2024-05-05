package service

import "errors"

type AccountService interface {
	CreateAccount(string, string) error
	GetUserInfo(string) (map[string]interface{}, error)
	PatchUser(string, string, string) error
	DeleteAccount(string)
}

type accountService struct {
	userData map[string]interface{}
}

func NewAccountService() AccountService {
	userData := map[string]interface{}{}
	return &accountService{
		userData: userData,
	}
}

func (a *accountService) CreateAccount(userId, password string) error {
	if a.userData[userId] == nil {
		a.userData[userId] = map[string]interface{}{
			"password": password,
			"nickname": userId,
			"user_id":  userId,
		}
		return nil
	} else {
		return errors.New("user Exists")
	}
}

func (a *accountService) GetUserInfo(userId string) (map[string]interface{}, error) {
	if a.userData[userId] != nil {
		return a.userData[userId].(map[string]interface{}), nil
	} else {
		return nil, errors.New("not found")
	}
}

func (a *accountService) PatchUser(userId, comment, nickname string) error {
	if a.userData[userId] != nil {
		newData := map[string]interface{}{
			"comment":  comment,
			"nickname": nickname,
			"password": a.userData[userId].(map[string]interface{})["password"].(string),
			"user_id":  userId,
		}

		a.userData[userId] = newData
		return nil
	} else {
		return errors.New("not found")
	}
}

func (a *accountService) DeleteAccount(userId string) {
	delete(a.userData, userId)
}
