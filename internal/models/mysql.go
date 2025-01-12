package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user model
type User struct {
	gorm.Model
	ID        int    `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// relationships
	Tweets []Tweet `gorm:"foreignKey:UserID"`
}

// Tweet represents the tweets model
type Tweet struct {
	gorm.Model
	ID          int    `gorm:"primaryKey"`
	UserID      int    `gorm:"not null"`                                                                      // Ensure this column cannot be null
	User        User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Foreign key relationship
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Followers represents the tweets model
type Followers struct {
	gorm.Model
	ID              int  `gorm:"primaryKey"`
	UserID          int  `gorm:"not null"` // Ensure this column cannot be null
	User            User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FollowingUserID int
	FollowingUsers  User `gorm:"foreignKey:FollowingUserID;references:ID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
