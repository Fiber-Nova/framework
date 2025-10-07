package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// MaintenanceNotice represents a maintenance message shown to users
// TableName: MaintenanceNotices
// Fields include scheduling window and optional scoping to station/instrument

type MaintenanceNotice struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Title         string     `gorm:"size:200;not null" json:"title"`
	Body          string     `json:"body"`
	Severity      string     `gorm:"type:ENUM('info','warning','critical');default:'info';index" json:"severity"`
	Status        string     `gorm:"type:ENUM('draft','published','archived');default:'published';index" json:"status"`
	Scope         string     `gorm:"type:ENUM('global','station','instrument');default:'global';index" json:"scope"`
	StationID     *uint      `gorm:"index" json:"station_id,omitempty"`
	InstrumentID  *uint      `gorm:"index" json:"instrument_id,omitempty"`
	EffectiveFrom *time.Time `gorm:"index" json:"effective_from,omitempty"`
	EffectiveTo   *time.Time `gorm:"index" json:"effective_to,omitempty"`
	Deleted       string     `gorm:"type:ENUM('Yes','No');default:'No'" json:"deleted"`
	CreatedByID   *uint      `gorm:"index" json:"created_by_id,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (MaintenanceNotice) TableName() string { return "MaintenanceNotices" }

// CreateMaintenanceNotice inserts a new notice
func CreateMaintenanceNotice(db *gorm.DB, n *MaintenanceNotice) error {
	if strings.TrimSpace(n.Title) == "" {
		return gorm.ErrInvalidData
	}
	return db.Create(n).Error
}

// UpdateMaintenanceNotice updates fields for a notice and returns the updated record
func UpdateMaintenanceNotice(db *gorm.DB, id uint, updates map[string]interface{}) (*MaintenanceNotice, error) {
	var cur MaintenanceNotice
	if err := db.First(&cur, id).Error; err != nil {
		return nil, err
	}
	if v, ok := updates["title"].(string); ok {
		updates["title"] = strings.TrimSpace(v)
	}
	if err := db.Model(&cur).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := db.First(&cur, id).Error; err != nil {
		return nil, err
	}
	return &cur, nil
}

// GetMaintenanceNoticeByID retrieves by ID
func GetMaintenanceNoticeByID(db *gorm.DB, id uint) (*MaintenanceNotice, error) {
	var n MaintenanceNotice
	if err := db.First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

// MaintenanceNoticeFilter supports filtering and pagination
type MaintenanceNoticeFilter struct {
	Severity       string
	Status         string
	Scope          string
	StationID      *uint
	InstrumentID   *uint
	ActiveOnly     bool
	IncludeDeleted bool
	Search         string     // title/body
	From           *time.Time // overlap window start
	To             *time.Time // overlap window end
	Page           int
	PageSize       int
}

// ListMaintenanceNotices lists notices according to filters and returns rows and total count
func ListMaintenanceNotices(db *gorm.DB, f MaintenanceNoticeFilter) ([]MaintenanceNotice, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize <= 0 || f.PageSize > 200 {
		f.PageSize = 20
	}
	q := db.Model(&MaintenanceNotice{})
	if f.Severity != "" {
		q = q.Where("severity = ?", f.Severity)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Scope != "" {
		q = q.Where("scope = ?", f.Scope)
	}
	if f.StationID != nil {
		q = q.Where("station_id = ?", *f.StationID)
	}
	if f.InstrumentID != nil {
		q = q.Where("instrument_id = ?", *f.InstrumentID)
	}
	if !f.IncludeDeleted {
		q = q.Where("deleted = 'No'")
	}
	if strings.TrimSpace(f.Search) != "" {
		like := "%" + strings.TrimSpace(f.Search) + "%"
		q = q.Where("title LIKE ? OR body LIKE ?", like, like)
	}
	// ActiveOnly: status=published and now within effective window (or open-ended)
	if f.ActiveOnly {
		now := time.Now()
		q = q.Where("status = 'published'")
		q = q.Where("(effective_from IS NULL OR effective_from <= ?)", now)
		q = q.Where("(effective_to IS NULL OR effective_to >= ?)", now)
	}
	// Optional overlap window filter
	if f.From != nil {
		q = q.Where("(effective_to IS NULL OR effective_to >= ?)", *f.From)
	}
	if f.To != nil {
		q = q.Where("(effective_from IS NULL OR effective_from <= ?)", *f.To)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []MaintenanceNotice
	if err := q.Order("effective_from DESC NULLS LAST, id DESC").Limit(f.PageSize).Offset((f.Page - 1) * f.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// MarkMaintenanceNoticeDeleted sets Deleted to "Yes"
func MarkMaintenanceNoticeDeleted(db *gorm.DB, id uint) error {
	return db.Model(&MaintenanceNotice{}).Where("id = ?", id).Update("deleted", "Yes").Error
}
