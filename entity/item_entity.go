package entity

import "time"

type Order struct {
	ID           int `gorm:"primaryKey"`
	CustomerName string
	OrderedAt    time.Time
	Items        []Item 
}

type Item struct {
	ID          int    `gorm:"primaryKey"`
	Code        string `gorm:"unique"`
	Description string
	Quantity    int
	OrderID     int
}
