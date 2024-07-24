package models

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Email     string
	FirstName string
	LastName  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
