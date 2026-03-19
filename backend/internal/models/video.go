package models

import (
	"time"
)

type Video struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Title             string    `json:"title" gorm:"not null"`
	Description       string    `json:"description"`
	YoutubeVideoID    string    `json:"youtube_video_id" gorm:"not null"`
	ThumbnailURL      string    `json:"thumbnail_url"`
	CustomThumbnailURL string   `json:"custom_thumbnail_url"`
	CategoryID        uint      `json:"category_id"`
	Category          Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Duration          string    `json:"duration"`
	ReleaseDate       string    `json:"release_date"`
	IsPremium         bool      `json:"is_premium" gorm:"default:false"`
	IsPublished       bool      `json:"is_published" gorm:"default:true"`
	ViewCount         int64     `json:"view_count" gorm:"default:0"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (v *Video) GetThumbnail() string {
	if v.CustomThumbnailURL != "" {
		return v.CustomThumbnailURL
	}
	if v.ThumbnailURL != "" {
		return v.ThumbnailURL
	}
	return "https://img.youtube.com/vi/" + v.YoutubeVideoID + "/maxresdefault.jpg"
}
