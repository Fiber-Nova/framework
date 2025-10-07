package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Configuration model for Configurations table
// TableName: Configurations
// Fields: ID, Key, Value, Description, UpdatedAt

type Configuration struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"unique;not null" json:"key"`
	Value       string    `gorm:"not null" json:"value"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Configuration) TableName() string {
	return "Configurations"
}

// GetConfigByKey retrieves a single configuration by its key
func GetConfigByKey(db *gorm.DB, key string) (*Configuration, error) {
	var cfg Configuration
	if err := db.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ListConfigs lists configurations, optionally filtered by a key prefix
func ListConfigs(db *gorm.DB, prefix string) ([]Configuration, error) {
	var items []Configuration
	q := db.Model(&Configuration{})
	if prefix != "" {
		q = q.Where("`key` LIKE ?", prefix+"%")
	}
	if err := q.Order("`key` asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// SetConfig updates an existing configuration by key or creates it if missing
func SetConfig(db *gorm.DB, key, value, description string) error {
	res := db.Model(&Configuration{}).Where("`key` = ?", key).Updates(map[string]interface{}{
		"value":       value,
		"description": description,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		cfg := Configuration{Key: key, Value: value, Description: description}
		return db.Create(&cfg).Error
	}
	return nil
}

// UpsertConfig inserts or updates a configuration using the key as the unique constraint
func UpsertConfig(db *gorm.DB, cfg *Configuration) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "description", "updated_at"}),
	}).Create(cfg).Error
}
