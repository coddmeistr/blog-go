package user

import (
	"github.com/maxik12233/blog/middleware"
	"github.com/maxik12233/blog/types"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req *CreateUserRequest) (uint, error)
	LoginUser(login string, pass string) (string, error)
	DeleteUser(id uint) error
	GetAll() ([]*User, error)
	GetOne(id uint) (*User, error)
	UpdateContactInfo(req *UpdateUserContactInfoRequest) error
	UpdatePersonalInfo(req *UpdateUserPersonalInfoRequest) error
	UpdateLocation(req *UpdateLocationRequest) error
}

type UserServiceImpl struct {
	repo   UserRepository
	logger *zap.Logger
}

func NewUserService(repo UserRepository, logger *zap.Logger) UserService {
	return &UserServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserServiceImpl) CreateUser(req *CreateUserRequest) (uint, error) {
	s.logger.Info("In CreateUser")

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		s.logger.Error("Error while hashing the password", zap.Error(err))
		return 0, ErrInternalError
	}
	req.Password = string(hash)

	user := User{
		Login:    req.Login,
		Password: req.Password,
		ContactInfo: &ContactInfo{
			Email: req.Email,
		},
		PersonalInfo: &PersonalInfo{
			Location: &Location{},
		},
		Role: []types.Role{
			{RoleDataID: 1},
			{RoleDataID: 2},
		},
	}

	id, err := s.repo.CreateUser(&user)
	if err != nil {
		s.logger.Error("Error while creating a user", zap.Error(err))
		return 0, err
	}

	return id, nil

}

func (s *UserServiceImpl) UpdateLocation(req *UpdateLocationRequest) error {
	s.logger.Info("In UpdateLocation")

	loc := Location{
		Country: req.Country,
		City:    req.City,
	}

	err := s.repo.UpdateLocation(req.ID, &loc)
	if err != nil {
		s.logger.Error("Error while updating location", zap.Error(err))
		return err
	}

	return nil
}

func (s *UserServiceImpl) UpdatePersonalInfo(req *UpdateUserPersonalInfoRequest) error {
	s.logger.Info("In UpdatePersonalInfo")

	pi := PersonalInfo{
		PersonalStatus: req.PersonalStatus,
		Description:    req.Description,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
	}

	err := s.repo.UpdateUserPersonalInfo(req.ID, &pi)
	if err != nil {
		s.logger.Error("Error while updating personal info", zap.Error(err))
		return err
	}

	return nil
}

func (s *UserServiceImpl) UpdateContactInfo(req *UpdateUserContactInfoRequest) error {
	s.logger.Info("In UpdateContactInfo")

	ci := ContactInfo{
		Email:  req.Email,
		Mobile: req.Mobile,
	}

	err := s.repo.UpdateUserContactInfo(req.ID, &ci)
	if err != nil {
		s.logger.Error("Error while updating contact info", zap.Error(err))
		return err
	}

	return nil
}

func (s *UserServiceImpl) LoginUser(login string, pass string) (string, error) {
	s.logger.Info("In LoginUser")

	var user *User
	var roles []uint
	user, err := s.repo.GetOneByLogin(login)
	if err != nil {
		s.logger.Error("Error while getting user by login", zap.Error(err))
		switch err {
		case ErrNotFound:
			return "", ErrInvalidLoginOrPassword
		default:
			return "", err
		}
	}

	roles, err = s.repo.GetUserRoles(user.ID)
	if err != nil {
		s.logger.Error("Error while getting user roles", zap.Error(err))
		return "", err
	}

	// Compare sent pass with hash pass
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		s.logger.Error("Error while comparing hash and password", zap.Error(err))
		return "", ErrInvalidLoginOrPassword
	}
	token, err := middleware.GenerateJWT(user.ID, roles)
	if err != nil {
		s.logger.Error("Error while creating jwt token", zap.Error(err))
		return "", ErrInternalError
	}

	return token, nil
}

func (s *UserServiceImpl) DeleteUser(id uint) error {
	s.logger.Info("In DeleteUser")

	err := s.repo.DeleteUser(id)
	if err != nil {
		s.logger.Error("Error while deleting a user", zap.Error(err))
		return err
	}

	return nil
}

func (s *UserServiceImpl) GetOne(id uint) (*User, error) {
	s.logger.Info("In GetOne")

	user, err := s.repo.GetOneById(id)
	if err != nil {
		s.logger.Error("Error while getting one user by id", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetAll() ([]*User, error) {
	s.logger.Info("In GetAll")

	users, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Error while getting all users", zap.Error(err))
		return nil, err
	}

	return users, nil
}
