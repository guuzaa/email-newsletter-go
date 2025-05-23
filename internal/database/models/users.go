package models

type User struct {
	ID       string `gorm:"column:user_id;not null;primaryKey;type:uuid"`
	Username string `gorm:"column:username;not null;unique"`
	Password string `gorm:"column:password_hash;not null"`
}
