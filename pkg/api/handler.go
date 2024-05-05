package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"practice/pkg/common"
	"practice/pkg/middleware"
	"practice/pkg/requests"
	"practice/pkg/responses"
	"practice/pkg/service"
	"practice/pkg/util"
	"strings"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(asvc service.AccountService) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/health", getHealthHandler())
	router.HandleFunc("/signup", getSignUpHandler(asvc))
	router.HandleFunc("/users/{user_id}", getUsersHandler(asvc))
	router.HandleFunc("/close", getCloseHandler(asvc))

	return router
}

func getSignUpHandler(accountService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "method not allowed", Message: "Method Not Allowed"}, 405)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "error parsing request body", Message: "Account creation failed"}, 400)
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

func getUsersHandler(accountService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" && req.Method != "PATCH" {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "method not allowed", Message: "Method Not Allowed"}, 405)
			return
		}

		vars := mux.Vars(req)
		if vars["user_id"] == "" {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "user_id is required", Message: "Account creation failed"}, 400)
			return
		}

		if util.SpacesOrControlCodes(vars["user_id"]) {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "user_id contains invalid characters", Message: "Account creation failed"}, 400)
			return
		}

		if !util.IsValidHalfWidthASCII(vars["user_id"], 6, 20) {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "user_id length or characters invalid", Message: "Account creation failed"}, 400)
			return
		}

		if req.Method == "GET" {
			performGetUserInfo(w, accountService, vars["user_id"])
		} else {
			performPatchUserInfo(w, req, accountService, vars["user_id"])
		}
	}
}

func performGetUserInfo(w http.ResponseWriter, accountService service.AccountService, userId string) {
	info, err := accountService.GetUserInfo(userId)

	if err != nil {
		middleware.WriteResponse(w, &responses.ErrorMessage{Message: "No user found"}, 404)
		return
	}

	if info["comment"] != nil && info["comment"].(string) != "" {
		middleware.WriteResponse(w, &responses.GetUserInfoConf{Message: "User details by user_id", User: common.UserInfo{UserId: userId, Nickname: info["nickname"].(string), Comment: info["comment"].(string)}}, 200)
		return
	}

	middleware.WriteResponse(w, &responses.GetUserInfo{Message: "User details by user_id", User: responses.User{UserId: userId, Nickname: info["nickname"].(string)}}, 200)
}

func performPatchUserInfo(w http.ResponseWriter, req *http.Request, accountService service.AccountService, userId string) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "error parsing request body", Message: "Patch user info failed"}, 400)
		return
	}

	var patchUserReq common.UserInfo
	err = json.Unmarshal(body, &patchUserReq)
	if err != nil {
		middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "error parsing request body", Message: "Patch user info failed"}, 400)
		return
	}

	if patchUserReq.Comment == "" && patchUserReq.Nickname == "" {
		middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "required comment or nickname", Message: "User updation failed"}, 400)
		return
	}

	if patchUserReq.UserId != "" || patchUserReq.Password != "" {
		middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "not updatable user_id and password", Message: "User updation failed"}, 400)
		return
	}

	if len(patchUserReq.Comment) > 100 {
		middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "comment length invalid", Message: "Patch request invalid"}, 400)
		return
	}

	if len(patchUserReq.Nickname) > 30 {
		middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "nickname length invalid", Message: "Patch request invalid"}, 400)
		return
	}

	if !util.ValidString(patchUserReq.Comment) {
		middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "comment contains invalid characters", Message: "User updation failed"}, 400)
		return
	}

	if !util.ValidString(patchUserReq.Nickname) {
		middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "nickname contains invalid characters", Message: "User updation failed"}, 400)
		return
	}

	if patchUserReq.Nickname == "" {
		patchUserReq.Nickname = userId
	}

	err = accountService.PatchUser(userId, patchUserReq.Comment, patchUserReq.Nickname)

	if err != nil {
		middleware.WriteResponse(w, &responses.ErrorMessage{Message: "No user found"}, 404)
		return
	}

	middleware.WriteResponse(w, &responses.PatchUserInfo{Message: "User successfully updated", Recipe: []responses.CommNick{{Comment: patchUserReq.Comment, Nickname: patchUserReq.Nickname}}}, 200)
}

func getCloseHandler(accountService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "method not allowed", Message: "Method Not Allowed"}, 405)
			return
		}

		header := req.Header.Get("Authorization")
		str, err := util.ParseAuthHeader(header)
		if err != nil {
			middleware.WriteResponse(w, &responses.ErrorMessage{Cause: "Unknown", Message: "Internal Server Error"}, 500)
			return
		}

		accountService.DeleteAccount(strings.Split(str, ":")[0])
		middleware.WriteResponse(w, &responses.DeleteAccount{Message: "Account and user successfully removed"}, 200)
	}
}

func getHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(200)
		if _, err := w.Write([]byte("UP")); err != nil {
			fmt.Print("Health Response Error")
		}
	}
}
