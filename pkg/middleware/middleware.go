package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"practice/pkg/responses"
	"practice/pkg/token"
	"strings"

	"github.com/gorilla/mux"
)

func WithAuth(verifier token.Verifier) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/health" || req.URL.Path == "/signup" {
				next.ServeHTTP(w, req)
				return
			}

			header, _ := GetAuthHeader(req)
			basic, _ := strings.CutPrefix(header, "Basic")
			trimmed := strings.TrimSpace(basic)

			vars := mux.Vars(req)
			var err error
			if vars["user_id"] != "" {
				err = verifier.Verify(trimmed, vars["user_id"])
			} else {
				err = verifier.Verify(trimmed, "")
			}

			if err != nil {
				if err.Error() == "Denied" {
					WriteResponse(w, map[string]string{"message": "No permission for update"}, 403)
					return
				}

				WriteResponse(w, map[string]string{"message": "Authentication failed"}, 401)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

func GetAuthHeader(req *http.Request) (string, error) {
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
