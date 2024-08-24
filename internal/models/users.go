package models

import (
	"time"
)

type User struct {
	ID            uint      `json:"user_id" gorm:"primaryKey"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Password      string    `json:"-"`
	Role          Role      `json:"role" gorm:"foreignKey:RoleID"`
	RoleID        uint      `json:"-"` // RoleID is needed for the foreign key relationship but is not exposed in JSON
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	VoteUpdatedAt time.Time `json:"vote_updated_at"`
	DeletedAt     time.Time `json:"-" gorm:"index"`
	Rating        int       `json:"rating"`
}
