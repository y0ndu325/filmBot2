package models

import "time"

type Movie struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"uniqueIndex;not null"`
	Overview    string `gorm:"type:text"`
	PosterPath  string `gorm:"type:varchar(255)"`
	ReleaseDate string `gorm:"type:varchar(10)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`	
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
