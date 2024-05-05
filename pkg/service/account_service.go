package service

import "errors"

type AccountService interface {
	CreateAccount(string, string) error
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
		}
		return nil
	} else {
		return errors.New("User Exists")
	}
}
