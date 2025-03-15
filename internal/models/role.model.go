package models

type RoleModel struct {
	BaseModel
	Name  string       `gorm:"unique;not null"`
	Users []*UserModel `gorm:"many2many:user_has_roles;"`
}

func (u *RoleModel) TableName() string {
	return "roles"
}
