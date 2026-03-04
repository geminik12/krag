package model

import "time"

type Document struct {
	ID        uint   `gorm:"primarykey"`
	UserID    string `gorm:"type:varchar(255);not null;index"`
	Filename  string `gorm:"type:varchar(255);not null"`
	Key       string `gorm:"type:varchar(255);not null"` // MinIO key
	Status    int    `gorm:"type:tinyint;default:0"`     // 0:Pending, 1:Processing, 2:Success, 3:Failed
	Error     string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	DocStatusPending    = 0
	DocStatusProcessing = 1
	DocStatusSuccess    = 2
	DocStatusFailed     = 3
)
