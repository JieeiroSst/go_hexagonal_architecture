package domain

import "time"

type User struct {
	ID           string `gorm:"primaryKey"`
	Name         string
	Email        string `gorm:"uniqueIndex"`
	Password     string
	LastActiveAt time.Time
	CreatedAt    time.Time
}
