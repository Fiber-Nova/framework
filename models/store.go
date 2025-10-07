package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// InventoryStore represents an instrument stored in inventory
// TableName: Stores
// Fields: ID, Location, Name, Code, CreatedAt, UpdatedAt, Deleted

type InventoryStore struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Location  string    `gorm:"size:200" json:"location"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Name      string    `gorm:"size:200;not null" json:"name"`
	Code      string    `gorm:"size:100;unique;not null;index" json:"code"`
	Deleted   string    `gorm:"type:ENUM('Yes','No');default:'No'" json:"deleted"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (InventoryStore) TableName() string { return "Stores" }

// CreateInventoryStore inserts a new store record
func CreateInventoryStore(db *gorm.DB, s *InventoryStore) error {
	if strings.TrimSpace(s.Name) == "" || strings.TrimSpace(s.Code) == "" {
		return gorm.ErrInvalidData
	}
	return db.Create(s).Error
}

// UpdateInventoryStore updates selected fields by ID
func UpdateInventoryStore(db *gorm.DB, id uint, updates map[string]interface{}) (*InventoryStore, error) {
	var cur InventoryStore
	if err := db.First(&cur, id).Error; err != nil {
		return nil, err
	}
	if v, ok := updates["name"].(string); ok {
		updates["name"] = strings.TrimSpace(v)
	}
	if v, ok := updates["code"].(string); ok {
		updates["code"] = strings.TrimSpace(v)
	}
	if err := db.Model(&cur).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := db.First(&cur, id).Error; err != nil {
		return nil, err
	}
	return &cur, nil
}

// GetInventoryStoreByID returns a single store row
func GetInventoryStoreByID(db *gorm.DB, id uint) (*InventoryStore, error) {
	var s InventoryStore
	if err := db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// ListInventoryStores returns stores filtered by optional search and pagination with optional sorting (name, code, location)
func ListInventoryStores(db *gorm.DB, search string, page, pageSize int, sortField, sortOrder string) ([]InventoryStore, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}
	q := db.Model(&InventoryStore{}).Where("deleted = 'No'")
	if s := strings.TrimSpace(search); s != "" {
		like := "%" + s + "%"
		q = q.Where("location LIKE ? OR name LIKE ? OR code LIKE ?", like, like, like)
	}
	// Safe sorting
	col := "name"
	switch strings.ToLower(strings.TrimSpace(sortField)) {
	case "code":
		col = "code"
	case "location":
		col = "location"
	case "name":
		fallthrough
	default:
		col = "name"
	}
	dir := "ASC"
	if strings.ToLower(strings.TrimSpace(sortOrder)) == "desc" {
		dir = "DESC"
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []InventoryStore
	if err := q.Order(fmt.Sprintf("%s %s", col, dir)).Limit(pageSize).Offset((page - 1) * pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// DeleteInventoryStore performs a hard delete by ID
func DeleteInventoryStore(db *gorm.DB, id uint) error { return db.Delete(&InventoryStore{}, id).Error }

// MarkInventoryStoreDeleted marks as deleted
func MarkInventoryStoreDeleted(db *gorm.DB, id uint) error {
	return db.Model(&InventoryStore{}).Where("id = ?", id).Update("deleted", "Yes").Error
}

// InventoryStorePage wraps a paginated result set for stores
type InventoryStorePage struct {
	Data       []InventoryStore `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
	HasNext    bool             `json:"has_next"`
	HasPrev    bool             `json:"has_prev"`
}

// LoadInventoryStorePage loads a page of stores with pagination metadata
func LoadInventoryStorePage(db *gorm.DB, search string, page, pageSize int, sortField, sortOrder string) (InventoryStorePage, error) {
	rows, total, err := ListInventoryStores(db, search, page, pageSize, sortField, sortOrder)
	if err != nil {
		return InventoryStorePage{}, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return InventoryStorePage{
		Data:       rows,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}, nil
}
