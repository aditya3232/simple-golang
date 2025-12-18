package model

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);unique;not null;index:idx_users_email"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Phone     string    `gorm:"type:varchar(17)"`
	Address   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt *time.Time
}

func (User) TableName() string {
	return "users"
}
