package models

import (
	"time"
)

type Message struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:varchar(120);not null"`
	Phone     string    `json:"phone" gorm:"type:varchar(15);not null"`
	Status    bool      `json:"status" gorm:"default:false"`
	MessageID string    `json:"message_id" gorm:"type:varchar(100)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
