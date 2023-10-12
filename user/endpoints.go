package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type UserEndpoints struct {
	CreateUser             endpoint.Endpoint
	DeleteUser             endpoint.Endpoint
	LoginUser              endpoint.Endpoint
	GetAllUsers            endpoint.Endpoint
	GetOneUser             endpoint.Endpoint
	UpdateUserContactInfo  endpoint.Endpoint
	UpdateUserPersonalInfo endpoint.Endpoint
	UpdateLocation         endpoint.Endpoint
}

func MakeUserEndpoints(s UserService) UserEndpoints {
	return UserEndpoints{
		CreateUser:             makeCreateUserEndpoint(s),
		DeleteUser:             makeDeleteUserEndpoint(s),
		LoginUser:              makeLoginUserEndpoint(s),
		GetAllUsers:            makeGetAllUsersEndpoint(s),
		GetOneUser:             makeGetOneUserEndpoint(s),
		UpdateUserContactInfo:  makeUpdateUserContactInfoEndpoint(s),
		UpdateUserPersonalInfo: makeUpdateUserPersonalInfoEndpoint(s),
		UpdateLocation:         makeUpdateLocationEndpoint(s),
	}
}

func makeUpdateLocationEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateLocationRequest)
		err := s.UpdateLocation(&req)
		if err != nil {
			return nil, err
		}
		return "Location info updated", nil
	}
}

func makeUpdateUserPersonalInfoEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserPersonalInfoRequest)
		err := s.UpdatePersonalInfo(&req)
		if err != nil {
			return nil, err
		}
		return "Personal info updated", nil
	}
}

func makeUpdateUserContactInfoEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserContactInfoRequest)
		err := s.UpdateContactInfo(&req)
		if err != nil {
			return nil, err
		}
		return "Contact info updated", nil
	}
}

func makeCreateUserEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		id, err := s.CreateUser(&req)
		if err != nil {
			return nil, err
		}
		return CreateUserResponse{ID: id, Message: "User created"}, nil
	}
}

func makeGetOneUserEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := s.GetOne(req.ID)
		if err != nil {
			return nil, err
		}
		return GetUserResponse{
			ID:           user.ID,
			Login:        user.Login,
			PersonalInfo: *user.PersonalInfo,
			ContactInfo:  *user.ContactInfo,
		}, nil
	}
}

func makeGetAllUsersEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.GetAll()
		if err != nil {
			return nil, err
		}
		var usersresp []GetUserResponse
		for _, val := range users {
			usersresp = append(usersresp, GetUserResponse{
				ID:           val.ID,
				Login:        val.Login,
				PersonalInfo: *val.PersonalInfo,
				ContactInfo:  *val.ContactInfo,
			})
		}
		return GetAllUsersResponse{
			Users: usersresp,
		}, nil
	}
}

func makeDeleteUserEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteUserRequest)
		err := s.DeleteUser(req.ID)
		if err != nil {
			return nil, err
		}
		return "User deleted", nil
	}
}

func makeLoginUserEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginUserRequest)
		token, err := s.LoginUser(req.Login, req.Password)
		if err != nil {
			return nil, err
		}
		return LoginUserResponse{
			Token: token,
		}, nil
	}
}
