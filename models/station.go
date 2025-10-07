package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Station model
type Station struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	Name          string          `gorm:"not null" json:"name"`
	Location      string          `json:"location"`
	Latitude      decimal.Decimal `gorm:"type:decimal(18,12)" json:"latitude"`
	Longitude     decimal.Decimal `gorm:"type:decimal(18,12)" json:"longitude"`
	Active        bool            `gorm:"default:true" json:"active"`
	StationTypeID uint            `gorm:"not null;index" json:"station_type_id"`
	StationType   StationType     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"station_type,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

// TableName sets the table name to 'STN' for GORM
func (Station) TableName() string {
	return "Stations"
}

// AddStation creates a new station
func AddStation(db *gorm.DB, station *Station) error {
	return db.Create(station).Error
}

// DeleteStation deletes a station by ID
func DeleteStation(db *gorm.DB, stationID uint) error {
	return db.Delete(&Station{}, stationID).Error
}

// DeactivateStation sets the station's Active field to false
func DeactivateStation(db *gorm.DB, stationID uint) error {
	return db.Model(&Station{}).Where("id = ?", stationID).Update("active", false).Error
}
