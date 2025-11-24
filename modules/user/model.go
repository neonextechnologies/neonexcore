package user

import (
	"time"

	"gorm.io/gorm"
)

// User model represents a user in the database
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Age       int            `gorm:"default:0" json:"age"`
	Active    bool           `gorm:"default:true" json:"active"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
