package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Seeder interface for database seeding
type Seeder interface {
	Seed(ctx context.Context, db *gorm.DB) error
}

// SeederManager manages database seeders
type SeederManager struct {
	db      *gorm.DB
	seeders []Seeder
}

// NewSeederManager creates a new seeder manager
func NewSeederManager(db *gorm.DB) *SeederManager {
	return &SeederManager{
		db:      db,
		seeders: make([]Seeder, 0),
	}
}

// Register registers a seeder
func (sm *SeederManager) Register(seeder Seeder) {
	sm.seeders = append(sm.seeders, seeder)
}

// Run runs all registered seeders
func (sm *SeederManager) Run(ctx context.Context) error {
	if len(sm.seeders) == 0 {
		fmt.Println("⚠️  No seeders registered")
		return nil
	}

	fmt.Printf("🌱 Running %d seeders...\n", len(sm.seeders))

	for _, seeder := range sm.seeders {
		if err := seeder.Seed(ctx, sm.db); err != nil {
			return fmt.Errorf("seeder failed: %w", err)
		}
	}

	fmt.Println("✅ Database seeding completed")
	return nil
}
