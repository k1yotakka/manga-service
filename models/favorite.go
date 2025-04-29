package models

type Favorite struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint `json:"user_id"`
	MangaID uint `json:"manga_id"`
}
