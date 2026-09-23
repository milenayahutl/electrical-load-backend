package ds

import (
	"database/sql"
	"time"
)

const (
	StatusDraft     = "черновик"
	StatusPublished = "опубликован"
	StatusDeleted   = "удален"
)

type ElectricityConsumer struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(255)"`
	Status      string `gorm:"type:varchar(20);not null"`
	Image       string `gorm:"type:varchar(255)"`
	Video       string `gorm:"type:varchar(255)"`

	PowerKW  float64 `gorm:"not null"`
	CurrentA float64 `gorm:"not null"`

	DateCreate    time.Time    `gorm:"not null"`
	CreatorID     uint         `gorm:"not null"`
	DateFormation sql.NullTime `gorm:"default:null"`

	Creator User   `gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT;"`
	Likes   []Like `gorm:"foreignKey:ElectricityConsumerID"`
}
