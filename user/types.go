package user

import (
	"github.com/maxik12233/blog/types"
)

/* types.go file stores golang structs for
   automigration and maintaining with GORM */

type Location struct {
	ID      uint   `gorm:"primaryKey"`
	Country string `gorm:"default:" json:"country"`
	City    string `gorm:"default:" json:"city"`
}

type ContactInfo struct {
	ID     uint   `gorm:"primaryKey"`
	Mobile string `gorm:"default:"`
	Email  string `gorm:"default:" json:"email"`
}

type PersonalInfo struct {
	ID             uint   `gorm:"primaryKey"`
	FirstName      string `gorm:"default:" json:"name"`
	LastName       string `gorm:"default:" json:"surname"`
	PersonalStatus string `gorm:"default:" json:"status"`
	Description    string `gorm:"default:" json:"descr"`

	LocationID uint      `gorm:"default:null" json:"-"`
	Location   *Location `gorm:"constraint:OnDelete:SET NULL; default:null"`
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Login    string `gorm:"unique" json:"login,omitempty"`
	Password string `json:"-,omitempty"`

	PersonalInfoID uint            `gorm:"default:null" json:"-"`
	PersonalInfo   *PersonalInfo   `gorm:"constraint:OnDelete:SET NULL; default:null" json:"personal,omitempty"`
	ContactInfoID  uint            `gorm:"default:null" json:"-"`
	ContactInfo    *ContactInfo    `gorm:"constraint:OnDelete:SET NULL; default:null" json:"contact,omitempty"`
	Role           []types.Role    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Article        []types.Article `gorm:"constraint:OnDelete:CASCADE;foreignKey:AuthorID" json:"-"`
	Comment        []types.Comment `gorm:"constraint:OnDelete:CASCADE;foreignKey:AuthorID" json:"-"`
	Like           []types.Like    `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID" json:"-"`
}
