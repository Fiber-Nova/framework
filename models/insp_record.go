package models

import (
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

// InspectionForm represents a submitted station inspection form
// TableName: InspectionForms
// Core fields plus a flexible JSON payload for dynamic form items

type InspectionRecord struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	StationID     uint            `gorm:"index:idx_insp_forms_station_id" json:"station_id"`
	Station       *Station        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	InstrumentID  *uint           `gorm:"index" json:"instrument_id,omitempty"`
	Instrument    *Instrument     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	SubmittedByID *uint           `gorm:"index:idx_insp_forms_submitter" json:"submitted_by_id,omitempty"`
	SubmittedBy   *User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	VisitDate     time.Time       `gorm:"index:idx_insp_forms_visit_date" json:"visit_date"`
	Status        string          `gorm:"type:ENUM('draft','submitted','approved','rejected');default:'submitted';index:idx_insp_forms_status" json:"status"`
	Title         string          `gorm:"size:200" json:"title"`
	Remarks       string          `json:"remarks"`
	Data          json.RawMessage `json:"data"`
	Deleted       string          `gorm:"type:ENUM('Yes','No');default:'No'" json:"deleted"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (InspectionRecord) TableName() string { return "InspectionForms" }

// InspectionFormAttachment stores metadata for files uploaded with an inspection form
// TableName: InspectionFormAttachments

type InspectionFormAttachment struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	InspectionFormID uint              `gorm:"index" json:"inspection_form_id"`
	InspectionForm   *InspectionRecord `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	FileName         string            `json:"file_name"`
	FilePath         string            `json:"file_path"`
	ContentType      string            `json:"content_type"`
	Size             int64             `json:"size"`
	UploadedAt       time.Time         `gorm:"autoCreateTime" json:"uploaded_at"`
}

func (InspectionFormAttachment) TableName() string { return "InspectionFormAttachments" }

// CreateInspectionForm inserts a new inspection form
func CreateInspectionForm(db *gorm.DB, f *InspectionRecord) error {
	return db.Create(f).Error
}

// UpdateInspectionForm updates a form by ID with the provided fields and returns the updated record
func UpdateInspectionForm(db *gorm.DB, id uint, updates map[string]interface{}) (*InspectionRecord, error) {
	var cur InspectionRecord
	if err := db.First(&cur, id).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&cur).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := db.First(&cur, id).Error; err != nil {
		return nil, err
	}
	return &cur, nil
}

// GetInspectionFormByID retrieves a single form by ID
func GetInspectionFormByID(db *gorm.DB, id uint) (*InspectionRecord, error) {
	var f InspectionRecord
	if err := db.First(&f, id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

// InspectionFormFilter defines optional filters and pagination for listing forms
type InspectionFormFilter struct {
	StationID      *uint
	SubmittedByID  *uint
	Status         string
	FromDate       *time.Time
	ToDate         *time.Time
	IncludeDeleted bool
	Search         string // matches title or remarks
	Page           int
	PageSize       int
}

// ListInspectionForms lists forms with filters and pagination, returning rows and total count
func ListInspectionForms(db *gorm.DB, f InspectionFormFilter) ([]InspectionRecord, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize <= 0 || f.PageSize > 200 {
		f.PageSize = 20
	}

	tx := db.Model(&InspectionRecord{})
	if f.StationID != nil {
		tx = tx.Where("station_id = ?", *f.StationID)
	}
	if f.SubmittedByID != nil {
		tx = tx.Where("submitted_by_id = ?", *f.SubmittedByID)
	}
	if f.Status != "" {
		tx = tx.Where("status = ?", f.Status)
	}
	if f.FromDate != nil {
		tx = tx.Where("visit_date >= ?", *f.FromDate)
	}
	if f.ToDate != nil {
		tx = tx.Where("visit_date <= ?", *f.ToDate)
	}
	if !f.IncludeDeleted {
		tx = tx.Where("deleted = 'No'")
	}
	if s := strings.TrimSpace(f.Search); s != "" {
		like := "%" + s + "%"
		tx = tx.Where("title LIKE ? OR remarks LIKE ?", like, like)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []InspectionRecord
	if err := tx.Order("visit_date DESC, id DESC").Limit(f.PageSize).Offset((f.Page - 1) * f.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// MarkInspectionFormDeleted marks a form as deleted (soft delete using the enum flag)
func MarkInspectionFormDeleted(db *gorm.DB, id uint) error {
	return db.Model(&InspectionRecord{}).Where("id = ?", id).Update("deleted", "Yes").Error
}

// AddInspectionAttachment links a file to a form
func AddInspectionAttachment(db *gorm.DB, att *InspectionFormAttachment) error {
	return db.Create(att).Error
}

// ListInspectionAttachments returns attachments for a form
func ListInspectionAttachments(db *gorm.DB, formID uint) ([]InspectionFormAttachment, error) {
	var items []InspectionFormAttachment
	if err := db.Where("inspection_form_id = ?", formID).Order("uploaded_at DESC, id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// DeleteInspectionAttachment removes an attachment by ID
func DeleteInspectionAttachment(db *gorm.DB, id uint) error {
	return db.Delete(&InspectionFormAttachment{}, id).Error
}
