package user

import (
	"github.com/maxik12233/blog/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(user *User) (uint, error)
	DeleteUser(id uint) error
	UpdateUserContactInfo(userid uint, contactinfo *ContactInfo) error
	UpdateUserPersonalInfo(userid uint, personalinfo *PersonalInfo) error
	UpdateLocation(userid uint, loc *Location) error
	GetAll() ([]*User, error)
	GetOneById(id uint) (*User, error)
	GetOneByLogin(login string) (*User, error)
	GetUserRoles(id uint) ([]uint, error)
}

type UserRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepo(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &UserRepo{
		db:     db,
		logger: logger,
	}
}

func (repo *UserRepo) UpdateLocation(userid uint, loc *Location) error {
	repo.logger.Info("In UpdateLocation")

	var personalid int
	result := repo.db.Table("users").Select("personal_info_id").Where("id = ?", userid).Limit(1).Scan(&personalid)
	if result.Error != nil {
		repo.logger.Error("Error while updating location", zap.Error(result.Error))
		return ErrInternalError
	}
	if result.RowsAffected == 0 {
		repo.logger.Error("Error while updating location", zap.Error(result.Error))
		return ErrInvalidId
	}

	var id int
	result = repo.db.Table("personal_infos").Select("location_id").Where("id = ?", personalid).Limit(1).Scan(&id)
	if result.Error != nil {
		repo.logger.Error("Error while updating location", zap.Error(result.Error))
		return ErrInternalError
	}
	if result.RowsAffected == 0 {
		repo.logger.Error("Error while updating location", zap.Error(result.Error))
		return ErrInternalError
	}

	loc.ID = uint(id)

	result = repo.db.Save(&loc)
	if result.Error != nil {
		repo.logger.Error("Error while updating location", zap.Error(result.Error))
		return ErrInternalError
	}

	return nil
}

func (repo *UserRepo) UpdateUserPersonalInfo(userid uint, personalinfo *PersonalInfo) error {
	repo.logger.Info("In UpdateUserPersonalInfo")

	var id int
	result := repo.db.Table("users").Select("personal_info_id").Where("id = ?", userid).Limit(1).Scan(&id)
	if result.Error != nil {
		repo.logger.Error("Error while updating user's personal info", zap.Error(result.Error))
		return ErrInternalError
	}
	if result.RowsAffected == 0 {
		repo.logger.Error("Error while updating user's personal info", zap.Error(result.Error))
		return ErrInvalidId
	}

	personalinfo.ID = uint(id)

	result = repo.db.Omit("location_id").Save(&personalinfo)
	if result.Error != nil {
		repo.logger.Error("Error while updating user's personal info", zap.Error(result.Error))
		return ErrInternalError
	}

	return nil
}

func (repo *UserRepo) UpdateUserContactInfo(userid uint, contactinfo *ContactInfo) error {
	repo.logger.Info("In UpdateUserContactInfo")

	var id int
	result := repo.db.Table("users").Select("contact_info_id").Where("id = ?", userid).Limit(1).Scan(&id)
	if result.Error != nil {
		repo.logger.Error("Error while updating user's contact info", zap.Error(result.Error))
		return ErrInternalError
	}
	if result.RowsAffected == 0 {
		repo.logger.Error("Error while updating user's conctact info", zap.Error(result.Error))
		return ErrInvalidId
	}

	contactinfo.ID = uint(id)

	result = repo.db.Save(&contactinfo)
	if result.Error != nil {
		repo.logger.Error("Error while updating user's contact info", zap.Error(result.Error))
		return ErrInternalError
	}

	return nil
}

func (repo *UserRepo) CreateUser(user *User) (uint, error) {
	repo.logger.Info("In CreateUser")

	if result := repo.db.Unscoped().Create(user); result.Error != nil {
		repo.logger.Error("Error while creating user", zap.Error(result.Error))
		return 0, ErrInternalError
	}
	return user.ID, nil
}

func (repo *UserRepo) GetUserRoles(id uint) ([]uint, error) {
	repo.logger.Info("In GetUserRoles")

	var roles []*types.Role
	if result := repo.db.Where("user_id = ?", id).Find(&roles); result.Error != nil {
		repo.logger.Error("Error while fetching user roles from db", zap.Error(result.Error))
		return nil, ErrInternalError
	}
	rolesids := make([]uint, 0)
	for _, val := range roles {
		rolesids = append(rolesids, val.RoleDataID)
	}
	return rolesids, nil
}

func (repo *UserRepo) DeleteUser(id uint) error {
	repo.logger.Info("In DeleteUser")

	var user *User
	result := repo.db.Find(&user, id)
	if result.Error != nil {
		repo.logger.Error("Error while deleting user from db", zap.Error(result.Error))
		return ErrInternalError
	}
	if result.RowsAffected == 0 {
		repo.logger.Info("Users not found by id while deleting user")
		return ErrNotFound
	}
	// Delete assosiations
	// TODO: Do it so you dont need to rewrite it after another new assotiation appears
	var personalinfo PersonalInfo
	repo.db.Find(&personalinfo, user.PersonalInfoID)
	repo.db.Unscoped().Where("ID = ?", personalinfo.LocationID).Delete(&Location{})
	repo.db.Unscoped().Where("ID = ?", user.PersonalInfoID).Delete(&PersonalInfo{})
	repo.db.Unscoped().Where("ID = ?", user.ContactInfoID).Delete(&ContactInfo{})
	repo.db.Unscoped().Model(&user).Association("Role").Unscoped().Clear()

	result = repo.db.Unscoped().Delete(&user)
	if result.Error != nil {
		repo.logger.Error("Error while deleting user from db", zap.Error(result.Error))
		return ErrInternalError
	}

	return nil
}

func (repo *UserRepo) GetOneById(id uint) (*User, error) {
	repo.logger.Info("In GetOneById")

	var user *User
	result := repo.db.Where("ID = ?", id).Preload("PersonalInfo.Location").Preload(clause.Associations).Find(&user)
	if result.Error != nil {
		repo.logger.Error("Error while fetching user by id from db", zap.Error(result.Error))
		return nil, ErrInternalError
	}
	if result.RowsAffected == 0 {
		repo.logger.Info("Users not found by id")
		return nil, ErrNotFound
	}

	return user, nil
}

func (repo *UserRepo) GetOneByLogin(login string) (*User, error) {
	repo.logger.Info("In GetOneByLogin")

	var user *User
	result := repo.db.Where("login = ?", login).Preload("PersonalInfo.Location").Preload(clause.Associations).Find(&user)
	if result.Error != nil {
		repo.logger.Error("Error while fetching user by login from db", zap.Error(result.Error))
		return nil, ErrInternalError
	}
	if result.RowsAffected == 0 {
		repo.logger.Info("Users not found by login")
		return nil, ErrNotFound
	}

	return user, nil
}

func (repo *UserRepo) GetAll() ([]*User, error) {
	repo.logger.Info("In GetAll")

	var users []*User

	result := repo.db.Preload("PersonalInfo.Location").Preload(clause.Associations).Find(&users)
	if result.Error != nil {
		repo.logger.Error("Error while fetching users from db", zap.Error(result.Error))
		return nil, ErrInternalError
	}
	if result.RowsAffected == 0 {
		repo.logger.Info("Users not found while fetching all users")
		return nil, ErrNotFound
	}

	return users, nil
}
