package models

type Manga struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
}
