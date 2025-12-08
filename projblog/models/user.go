package models

type Users struct {
	Id       int    `jsn:"id" gorm:"id"`
	Username string `gorm:"unque;not null" json:"username"`
	Email    string `gorm:"unque;not null"json:"email"`
	Password string `gorm:"not null" json:"password"`
}

// 强制覆盖命名规则 绑定到users表 而不是userss表
func (u *Users) TableName() string {
	return "users"
}
