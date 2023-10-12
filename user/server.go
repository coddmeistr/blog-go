package user

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/maxik12233/blog/middleware"
	"github.com/maxik12233/blog/types"
)

func CreateNewServer(rg *mux.Router, endpoints UserEndpoints) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	usergroupAuth := rg.PathPrefix("/user").Subrouter()
	usergroupAuth.Use(middleware.LoggingMiddleware)
	usergroupAuth.Use(middleware.ValidateRolesMiddleware(
		[]uint{
			uint(types.RoleCommon),
		},
	))
	usergroupNoAuth := rg.PathPrefix("/user").Subrouter()

	usergroupAuth.Methods("PUT").Path("/{id}/personal/location").Handler(httptransport.NewServer(
		endpoints.UpdateLocation,
		decodeUpdateLocationRequest,
		encodeResponse,
		options...,
	))

	usergroupAuth.Methods("PUT").Path("/{id}/personal").Handler(httptransport.NewServer(
		endpoints.UpdateUserPersonalInfo,
		decodeUpdateUserPersonalInfoRequest,
		encodeResponse,
		options...,
	))

	usergroupAuth.Methods("PUT").Path("/{id}/contact").Handler(httptransport.NewServer(
		endpoints.UpdateUserContactInfo,
		decodeUpdateUserContactInfoRequest,
		encodeResponse,
		options...,
	))

	usergroupNoAuth.Methods("POST").Path("/").Handler(httptransport.NewServer(
		endpoints.CreateUser,
		decodeCreateUserRequest,
		encodeResponse,
		options...,
	))

	usergroupAuth.Methods("DELETE").Path("/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteUser,
		decodeDeleteUserRequest,
		encodeResponse,
		options...,
	))

	usergroupNoAuth.Methods("PUT").Path("/").Handler(httptransport.NewServer(
		endpoints.LoginUser,
		decodeLoginUserRequest,
		encodeLoginUserResponse,
		options...,
	))

	usergroupAuth.Methods("GET").Path("/{id}").Handler(httptransport.NewServer(
		endpoints.GetOneUser,
		decodeGetOneUserRequest,
		encodeResponse,
		options...,
	))

	usergroupAuth.Methods("GET").Path("/").Handler(httptransport.NewServer(
		endpoints.GetAllUsers,
		decodeGetAllUsersRequest,
		encodeResponse,
		options...,
	))

}
