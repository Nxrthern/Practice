package util

import (
	"encoding/base64"
	"errors"
	"strings"
)

func ParseAuthHeader(header string) (string, error) {
	basic, _ := strings.CutPrefix(header, "Basic")
	trimmed := strings.TrimSpace(basic)

	apiKey, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return "", errors.New("Unauthorized")
	} else {
		return string(apiKey), nil
	}
}
