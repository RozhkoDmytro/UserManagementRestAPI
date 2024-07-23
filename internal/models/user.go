package models

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
