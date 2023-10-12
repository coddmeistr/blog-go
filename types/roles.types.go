package types

type Roles uint

const (
	RoleGuest Roles = iota
	RoleCommon
	RoleModerator
	RoleAdmin
)

type RoleData struct {
	ID       uint `gorm:"primaryKey"`
	RoleName string
}

type Role struct {
	ID uint `gorm:"primaryKey"`

	UserID     uint
	RoleDataID uint
	RoleData   RoleData
}
