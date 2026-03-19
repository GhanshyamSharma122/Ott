package models

import (
	"time"
)

type WatchHistory struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id" gorm:"index;not null"`
	User            User      `json:"-" gorm:"foreignKey:UserID"`
	VideoID         uint      `json:"video_id" gorm:"not null"`
	Video           Video     `json:"video" gorm:"foreignKey:VideoID"`
	WatchedDuration float64   `json:"watched_duration"`
	TotalDuration   float64   `json:"total_duration"`
	LastWatchedAt   time.Time `json:"last_watched_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (w *WatchHistory) ProgressPercent() float64 {
	if w.TotalDuration == 0 {
		return 0
	}
	return (w.WatchedDuration / w.TotalDuration) * 100
}
