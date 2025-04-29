package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	MangaID   uint      `json:"manga_id"`
	UserID    uint      `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
