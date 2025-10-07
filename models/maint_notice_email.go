package models

import "time"

type MaintNoticeEmail struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Station      string     `gorm:"size:50;index" json:"station"`
	FromTime     *time.Time `json:"from_time"`
	ToTime       *time.Time `json:"to_time"`
	UntilFurther bool       `gorm:"default:false" json:"until_further"`
	To           string     `gorm:"size:255" json:"to"`
	Template     string     `gorm:"size:150" json:"template"`
	Subject      string     `gorm:"size:255" json:"subject"`
	Body         string     `gorm:"type:longtext" json:"body"`
	SentBy       string     `gorm:"size:100;index" json:"sent_by"`
	SentAt       time.Time  `json:"sent_at"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (MaintNoticeEmail) TableName() string { return "MaintNoticeEmails" }
