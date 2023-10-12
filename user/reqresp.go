package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/maxik12233/blog/common"
)

/* reqresp.go file stores golang structs, types and functions
   to perform request response logic */

var ErrBadRequest = errors.New("Bad request")
var ErrInvalidId = errors.New("Invalid id")
var ErrNotFound = errors.New("Not found")
var ErrInternalError = errors.New("Internal error")
var ErrInvalidLoginOrPassword = errors.New("Invalid login or password")

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	case ErrInvalidId, ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case ErrInvalidLoginOrPassword:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

type UpdateLocationRequest struct {
	ID uint `json:"id"`
	Location
}

type UpdateUserPersonalInfoRequest struct {
	ID uint `json:"id"`
	PersonalInfo
}

type UpdateUserContactInfoRequest struct {
	ID uint `json:"id"`
	ContactInfo
}

type CreateUserRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

type GetUserRequest struct {
	ID uint `json:"id"`
}

type GetUserResponse struct {
	ID           uint         `json:"id"`
	Login        string       `json:"login"`
	PersonalInfo PersonalInfo `json:"personal"`
	ContactInfo  ContactInfo  `json:"contact"`
}

type GetAllUsersResponse struct {
	Users []GetUserResponse
}

type DeleteUserRequest struct {
	ID uint `json:"id"`
}

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

func decodeUpdateLocationRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req UpdateLocationRequest
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return nil, ErrInvalidId
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}
	req.ID = uint(id)
	return req, err
}

func decodeUpdateUserPersonalInfoRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req UpdateUserPersonalInfoRequest
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return nil, ErrInvalidId
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}
	req.ID = uint(id)
	return req, err
}

func decodeUpdateUserContactInfoRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req UpdateUserContactInfoRequest
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return nil, ErrInvalidId
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}
	req.ID = uint(id)
	return req, err
}

func decodeCreateUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}
	return req, err
}

func decodeDeleteUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return nil, ErrBadRequest
	}
	return DeleteUserRequest{
		ID: uint(id),
	}, nil
}

func decodeLoginUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}
	return LoginUserRequest{
		Login:    req.Login,
		Password: req.Password,
	}, nil
}

func decodeGetOneUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return nil, ErrBadRequest
	}
	return GetUserRequest{
		ID: uint(id),
	}, nil
}

func decodeGetAllUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeLoginUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {

	resp := response.(LoginUserResponse)
	// set cookie for storing token
	cookie := http.Cookie{}
	cookie.Name = common.JWT_TOKEN_NAME
	cookie.Value = resp.Token
	cookie.Expires = time.Now().Add(time.Hour * time.Duration(common.JWT_TOKEN_EXP_HOURS))
	cookie.Secure = false
	cookie.HttpOnly = true
	cookie.Path = "/"
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
