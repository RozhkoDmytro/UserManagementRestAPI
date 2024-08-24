package models

import (
	"time"

	"gorm.io/gorm"
)

type Vote struct {
	ID        uint      `json:"vote_id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`    // Voting user ID
	ProfileID uint      `json:"profile_id"` // ID of the profile being voted for
	Value     int       `json:"value"`      // Voice value (+1 or -1)
	CreatedAt time.Time `json:"created_at"` // Voting time
}

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
