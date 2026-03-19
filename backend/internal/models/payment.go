package models

import (
	"time"
)

type Payment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency" gorm:"default:INR"`
	Status    string    `json:"status" gorm:"default:success"`
	CreatedAt time.Time `json:"created_at"`
}
