package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string `gorm:"primaryKey; type:varchar(36)"`
	NIP       string `gorm:"type:varchar(18);uniqueIndex"`
	NIK       string `gorm:"type:varchar(16);uniqueIndex"`
	Username  string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password  string `gorm:"type:varchar(255);not null"`
	Role      int
	Position  string
	Name      string
	Telp      string
	Sex       string `gorm:"type:varchar(1)"`
	Address   string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Users []User
