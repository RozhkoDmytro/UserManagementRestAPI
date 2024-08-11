package models

import "time"

type Role struct {
	ID   uint   `json:"role_id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
}

type User struct {
	ID        uint      `json:"user_id" gorm:"primaryKey"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"-"`
	Role      Role      `json:"role" gorm:"foreignKey:RoleID"`
	RoleID    uint      `json:"-"` // RoleID is needed for the foreign key relationship but is not exposed in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-" gorm:"index"`
}

const (
	StrAdmin     = "admin"
	StrModerator = "moderator"
	StrUser      = "user"
)

// Define a custom type for the context key
type contextKey string

const (
	RoleContextKey  contextKey = "role"
	EmailContextKey contextKey = "email"
)
