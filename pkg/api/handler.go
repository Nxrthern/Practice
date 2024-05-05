package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"practice/pkg/middleware"
	"practice/pkg/requests"
	"practice/pkg/responses"
	"practice/pkg/service"
	"practice/pkg/util"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(hsvc service.HealthService, asvc service.AccountService) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/health", getHealthHandler(hsvc))
	router.HandleFunc("/signup", getSignUpHandler(asvc))
	router.HandleFunc("/health", getHealthHandler(hsvc))
	router.HandleFunc("/health", getHealthHandler(hsvc))
	router.HandleFunc("/health", getHealthHandler(hsvc))

	return router
}

func getSignUpHandler(accountService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var signUpReq requests.SignUpRequest
		err = json.Unmarshal(body, &signUpReq)
		if err != nil {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "error parsing request body", Message: "Account creation failed"}, 400)
			return
		}

		if signUpReq.UserId == "" {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "user_id is required", Message: "Account creation failed"}, 400)
			return
		}

		if signUpReq.Password == "" {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "password is required", Message: "Account creation failed"}, 400)
			return
		}

		if util.SpacesOrControlCodes(signUpReq.UserId) {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "user_id contains invalid characters", Message: "Account creation failed"}, 400)
			return
		}

		if !util.IsValidHalfWidthASCII(signUpReq.UserId, 6, 20) {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "user_id length or characters invalid", Message: "Account creation failed"}, 400)
			return
		}

		if util.SpacesOrControlCodes(signUpReq.Password) {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "password contains invalid characters", Message: "Account creation failed"}, 400)
			return
		}

		if !util.IsValidHalfWidthASCII(signUpReq.Password, 8, 20) {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "password length or characters invalid", Message: "Account creation failed"}, 400)
			return
		}

		err = accountService.CreateAccount(signUpReq.UserId, signUpReq.Password)

		if err != nil {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "already same user_id is used", Message: "Account creation failed"}, 400)
			return
		}

		middleware.WriteResponse(w, &responses.SignUpResponse{User: responses.User{UserId: signUpReq.UserId, Nickname: signUpReq.UserId}, Message: "Account successfully created"}, 200)
	}
}

func getHealthHandler(svc service.HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(200)
		if _, err := w.Write([]byte("UP")); err != nil {
			fmt.Print("Health Response Error")
		}
	}
}
