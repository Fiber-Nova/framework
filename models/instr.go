package models

import (
	"time"

	"gorm.io/gorm"
)

// Instrument model for scientific/measurement instruments
// TableName: Instruments
// Fields: ID, Name, Type, SerialNumber, Location, Status, CreatedAt, UpdatedAt

type Instrument struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	Name           string          `gorm:"not null" json:"name"`
	Type           uint            `gorm:"not null;index" json:"type"` // references InstrumentType.ID
	InstrumentType InstrumentType  `gorm:"foreignKey:Type;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"instrument_type,omitempty"`
	SerialNumber   string          `gorm:"unique;not null" json:"serial_number"`
	Location       string          `json:"location"`
	Status         string          `json:"status"` // e.g. active, inactive, maintenance
	StoreID        *uint           `gorm:"index" json:"store_id,omitempty"`
	Store          *InventoryStore `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"store,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Instrument) TableName() string {
	return "Instruments"
}

// AddInstrument creates a new instrument
func AddInstrument(db *gorm.DB, instr *Instrument) error {
	return db.Create(instr).Error
}

// DeleteInstrument deletes an instrument by ID
func DeleteInstrument(db *gorm.DB, instrID uint) error {
	return db.Delete(&Instrument{}, instrID).Error
}

// UpdateInstrument updates an instrument's details
func UpdateInstrument(db *gorm.DB, instr *Instrument) error {
	return db.Save(instr).Error
}

// FindInstrumentByID retrieves an instrument by ID
func FindInstrumentByID(db *gorm.DB, instrID uint) (*Instrument, error) {
	var instr Instrument
	if err := db.First(&instr, instrID).Error; err != nil {
		return nil, err
	}
	return &instr, nil
}
