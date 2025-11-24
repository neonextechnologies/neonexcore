package user

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// UserSeeder seeds initial user data
type UserSeeder struct{}

// Seed implements the Seeder interface
func (s *UserSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	// Check if users already exist
	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		fmt.Println("⏭️  Users already seeded, skipping...")
		return nil
	}

	users := []User{
		{
			Name:     "Alice Johnson",
			Email:    "alice@example.com",
			Password: "hashed_password_here",
			Age:      28,
			Active:   true,
		},
		{
			Name:     "Bob Smith",
			Email:    "bob@example.com",
			Password: "hashed_password_here",
			Age:      35,
			Active:   true,
		},
		{
			Name:     "Charlie Brown",
			Email:    "charlie@example.com",
			Password: "hashed_password_here",
			Age:      42,
			Active:   true,
		},
		{
			Name:     "Diana Prince",
			Email:    "diana@example.com",
			Password: "hashed_password_here",
			Age:      30,
			Active:   false,
		},
	}

	result := db.Create(&users)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("✅ Seeded %d users\n", len(users))
	return nil
}
