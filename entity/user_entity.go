package entity

import "time"

type User struct {
	ID        int `gorm:"primaryKey"`
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Order     []Order
}
