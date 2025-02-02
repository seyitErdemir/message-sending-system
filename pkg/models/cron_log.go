package models

import (
	"time"
)

type CronLog struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Operation     string    `json:"operation" gorm:"type:varchar(50);not null"` // START, STOP, UPDATE
	MessageIDs    string    `json:"message_ids" gorm:"type:text"`               // Comma separated message IDs
	MessagesCount int       `json:"messages_count"`
	Status        bool      `json:"status"` // Success or Failure
	Description   string    `json:"description" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}
