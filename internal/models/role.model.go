package models

type RoleModel struct {
	BaseModel
	Name  string       `gorm:"unique;not null" json:"name"`
	Users []*UserModel `gorm:"many2many:user_has_roles;joinForeignKey:RoleID;joinReferences:UserID" json:"users"`
}

func (u *RoleModel) TableName() string {
	return "roles"
}
