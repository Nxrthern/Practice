package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"practice/pkg/responses"
	"practice/pkg/token"
	"regexp"

	"github.com/gorilla/mux"
)

var rx = regexp.MustCompile(`(?i)^Bearer (.*)?$`)

func WithAuth(verifier token.Verifier) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/health" {
				next.ServeHTTP(w, req)
				return
			}

			header, err := getAuthHeader(req)
			ValidateSecretKey(&header)
			if err != nil {
				http.Error(w, "Unauthorized", 401)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

func ValidateSecretKey(token *string) error {
	apiKey, err := base64.StdEncoding.DecodeString(*token)
	if err != nil {
		return errors.New("Unauthorized")
	}

	//TODO
	if apiKey != nil {
		return nil
	} else {
		return errors.New("Unauthorized")
	}
}

func getAuthHeader(req *http.Request) (string, error) {
	return req.Header.Get("Authorization"), nil
}

func WriteResponse(w http.ResponseWriter, response any, statusCode int) {
	jsonData, err := json.Marshal(response)
	if err != nil {
		WriteResponse(w, &responses.ErrorMessage{Cause: "Unknown", Message: "Internal Server Error"}, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonData); err != nil {
		fmt.Print("Response Error")
	}
}
