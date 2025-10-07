package models

import (
	"time"

	"gorm.io/gorm"
)

// StationType represents a type/category of station
// Table: StationTypes
// Fields: ID, Code, Name, Description, CreatedAt, UpdatedAt

type StationType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StationType) TableName() string { return "StationTypes" }

// CreateStationType inserts a new StationType
func CreateStationType(db *gorm.DB, st *StationType) error {
	return db.Create(st).Error
}

// UpdateStationType updates fields of a StationType by id
func UpdateStationType(db *gorm.DB, id uint, changes map[string]interface{}) error {
	if len(changes) == 0 {
		return nil
	}
	return db.Model(&StationType{}).Where("id = ?", id).Updates(changes).Error
}

// DeleteStationType removes a StationType by id
func DeleteStationType(db *gorm.DB, id uint) error {
	return db.Delete(&StationType{}, id).Error
}
