package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID   uint   `json:"role_id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
}

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

type Vote struct {
	ID        uint      `json:"vote_id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`    // Voting user ID
	ProfileID uint      `json:"profile_id"` // ID of the profile being voted for
	Value     int       `json:"value"`      // Voice value (+1 or -1)
	CreatedAt time.Time `json:"created_at"` // Voting time
}

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

// AfterSave - a hook to automatically update the rating after saving a vote
func (v *Vote) AfterSave(tx *gorm.DB) (err error) {
	var rating int
	// Calculate a new rating for the profile that was voted for
	err = tx.Model(&Vote{}).
		Where("profile_id = ?", v.ProfileID).
		Select("COALESCE(SUM(value), 0)").
		Scan(&rating).Error
	if err != nil {
		return err
	}

	// Update the rating in the user table
	err = tx.Model(&User{}).Where("id = ?", v.ProfileID).Update("rating", rating).Error
	if err != nil {
		return err
	}

	// Update the UpdatedAt field for the profile
	err = tx.Model(&User{}).Where("id = ?", v.UserID).Update("vote_updated_at", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}
