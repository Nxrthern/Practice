package token

import (
	"encoding/base64"
	"errors"
	"practice/pkg/service"
	"strings"
)

var _ Verifier = &jwtTokenVerifyer{}

const ()

var ()

func NewJwTVerifier(asvc service.AccountService) Verifier {
	return &jwtTokenVerifyer{
		asvc: asvc,
	}
}

type jwtTokenVerifyer struct {
	asvc service.AccountService
}

func (s *jwtTokenVerifyer) Verify(str string, modUser string) error {
	apiKey, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return errors.New("Unauthorized")
	}

	apiKeyParts := strings.Split(string(apiKey), ":")
	info, err := s.asvc.GetUserInfo(apiKeyParts[0])
	if err != nil {
		return errors.New("Unauthorized")
	}

	if modUser != "" && apiKeyParts[0] != modUser {
		return errors.New("Denied")
	}

	if info["password"] == apiKeyParts[1] {
		return nil
	} else {
		return errors.New("Unauthorized")
	}
}
