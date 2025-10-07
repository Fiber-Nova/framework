package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// InstrumentType represents a type/category of instrument
// TableName: InstrumentTypes
// Fields: ID, Code, Name, Description, CreatedAt, UpdatedAt

type InstrumentType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:50;unique;not null" json:"code"`
	Name        string    `gorm:"not null;unique" json:"name"`
	Category    string    `json:"category,omitempty"`
	Status      string    `json:"status,omitempty"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (InstrumentType) TableName() string { return "InstrumentTypes" }

// CreateInstrumentType inserts a new row, deriving a unique normalized Code when missing.
func CreateInstrumentType(db *gorm.DB, it *InstrumentType) error {
	if strings.TrimSpace(it.Name) == "" {
		return fmt.Errorf("name is required")
	}
	candidate := it.Code
	if strings.TrimSpace(candidate) == "" {
		candidate = it.Name
	}
	unique, err := ensureUniqueInstrTypeCode(db, candidate, 0)
	if err != nil {
		return err
	}
	it.Code = unique
	return db.Create(it).Error
}

// UpdateInstrumentType updates allowed fields and enforces code uniqueness.
func UpdateInstrumentType(db *gorm.DB, id uint, updates map[string]interface{}) (*InstrumentType, error) {
	var current InstrumentType
	if err := db.First(&current, id).Error; err != nil {
		return nil, err
	}

	if v, ok := updates["name"].(string); ok {
		updates["name"] = strings.TrimSpace(v)
	}
	if v, ok := updates["description"].(string); ok {
		updates["description"] = strings.TrimSpace(v)
	}

	if v, ok := updates["code"].(string); ok {
		candidate := v
		if strings.TrimSpace(candidate) == "" {
			if nameNew, ok2 := updates["name"].(string); ok2 && strings.TrimSpace(nameNew) != "" {
				candidate = nameNew
			} else {
				candidate = current.Name
			}
		}
		unique, err := ensureUniqueInstrTypeCode(db, candidate, id)
		if err != nil {
			return nil, err
		}
		updates["code"] = unique
	}

	if err := db.Model(&current).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := db.First(&current, id).Error; err != nil {
		return nil, err
	}
	return &current, nil
}

// DeleteInstrumentType performs a hard delete.
func DeleteInstrumentType(db *gorm.DB, id uint) error { return db.Delete(&InstrumentType{}, id).Error }

// GetInstrumentTypeByID fetches a single record by ID.
func GetInstrumentTypeByID(db *gorm.DB, id uint) (*InstrumentType, error) {
	var it InstrumentType
	if err := db.First(&it, id).Error; err != nil {
		return nil, err
	}
	return &it, nil
}

// GetInstrumentTypeByCode fetches by unique code.
func GetInstrumentTypeByCode(db *gorm.DB, code string) (*InstrumentType, error) {
	var it InstrumentType
	if err := db.Where("code = ?", code).First(&it).Error; err != nil {
		return nil, err
	}
	return &it, nil
}

// ListInstrumentTypes supports search (code/name) and pagination; page starts at 1.
func ListInstrumentTypes(db *gorm.DB, q string, page, pageSize int) ([]InstrumentType, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}
	var items []InstrumentType
	var total int64
	tx := db.Model(&InstrumentType{})
	if s := strings.TrimSpace(q); s != "" {
		like := "%" + s + "%"
		tx = tx.Where("code LIKE ? OR name LIKE ?", like, like)
	}
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("name ASC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// UpsertInstrumentTypeByCode inserts or updates matched by code.
func UpsertInstrumentTypeByCode(db *gorm.DB, it *InstrumentType) error {
	candidate := it.Code
	if strings.TrimSpace(candidate) == "" {
		candidate = it.Name
	}
	it.Code = normalizeInstrTypeCode(candidate)
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "description", "updated_at"}),
	}).Create(it).Error
}

// BackfillEmptyInstrTypeCodes generates unique codes for rows with empty/NULL code.
func BackfillEmptyInstrTypeCodes(db *gorm.DB) (int64, error) {
	var list []InstrumentType
	if err := db.Where("code = '' OR code IS NULL").Find(&list).Error; err != nil {
		return 0, err
	}
	var updated int64
	for i := range list {
		row := &list[i]
		candidate := row.Name
		if strings.TrimSpace(candidate) == "" {
			candidate = fmt.Sprintf("IT-%d", row.ID)
		}
		unique, err := ensureUniqueInstrTypeCode(db, candidate, row.ID)
		if err != nil {
			return updated, err
		}
		if err := db.Model(row).Update("code", unique).Error; err != nil {
			return updated, err
		}
		updated++
	}
	return updated, nil
}

// normalizeInstrTypeCode converts input to uppercase hyphenated code within 50 chars.
func normalizeInstrTypeCode(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	s = strings.Join(strings.Fields(s), "-")
	s = strings.ToUpper(s)
	if len(s) > 50 {
		s = s[:50]
	}
	return s
}

// ensureUniqueInstrTypeCode ensures code uniqueness, excluding a given ID.
func ensureUniqueInstrTypeCode(db *gorm.DB, candidate string, excludeID uint) (string, error) {
	base := normalizeInstrTypeCode(candidate)
	if base == "" {
		return "", fmt.Errorf("code or name required to generate code")
	}
	code := base
	i := 1
	for {
		var count int64
		q := db.Model(&InstrumentType{}).Where("code = ?", code)
		if excludeID > 0 {
			q = q.Where("id <> ?", excludeID)
		}
		if err := q.Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
		suffix := fmt.Sprintf("-%d", i)
		trunc := base
		if len(trunc)+len(suffix) > 50 {
			trunc = trunc[:50-len(suffix)]
		}
		code = trunc + suffix
		i++
		if i > 1000 {
			return "", fmt.Errorf("unable to generate unique code")
		}
	}
}
