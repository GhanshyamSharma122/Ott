package models

import (
	"time"
)

type WatchlistItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	VideoID   uint      `json:"video_id" gorm:"not null"`
	Video     Video     `json:"video" gorm:"foreignKey:VideoID"`
	CreatedAt time.Time `json:"created_at"`
}
