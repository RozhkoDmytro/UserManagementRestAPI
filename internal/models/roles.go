package models

// Define a custom type for the context key
type contextKey string

const (
	StrAdmin     = "admin"
	StrModerator = "moderator"
	StrUser      = "user"
)

const (
	RoleContextKey  contextKey = "role"
	EmailContextKey contextKey = "email"
	IDContextKey    contextKey = "id"
)

type Role struct {
	ID   uint   `json:"role_id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
}
