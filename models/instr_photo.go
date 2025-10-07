package models

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

// InstrumentPhoto stores an instrument's photo binary in the database
// Allowed formats: PNG, JPEG, WebP
// TableName: InstrumentPhotos

type InstrumentPhoto struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	InstrumentID uint      `gorm:"not null;index" json:"instrument_id"`
	Filename     string    `gorm:"size:255;not null" json:"filename"`
	ContentType  string    `gorm:"size:100;not null" json:"content_type"`
	Size         int64     `gorm:"not null" json:"size"`
	Data         []byte    `gorm:"type:longblob;not null" json:"-"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (InstrumentPhoto) TableName() string { return "InstrumentPhotos" }

func (p *InstrumentPhoto) BeforeCreate(tx *gorm.DB) error { return p.validate() }
func (p *InstrumentPhoto) BeforeUpdate(tx *gorm.DB) error { return p.validate() }

func (p *InstrumentPhoto) validate() error {
	if p.InstrumentID == 0 {
		return errors.New("instrument id is required")
	}
	if len(p.Data) == 0 {
		return errors.New("image data is required")
	}
	if p.ContentType == "" {
		p.ContentType = http.DetectContentType(p.Data)
	}
	// reuse allowedPhotoMIMEs from stn_photo.go within models package
	if _, ok := allowedPhotoMIMEs[p.ContentType]; !ok {
		return errors.New("unsupported image type; allowed: PNG, JPEG, WebP")
	}
	// Normalize filename extension to match content type
	ext := strings.ToLower(filepath.Ext(p.Filename))
	switch p.ContentType {
	case "image/png":
		if ext != ".png" {
			p.Filename = strings.TrimSuffix(p.Filename, ext) + ".png"
		}
	case "image/jpeg":
		if ext != ".jpg" && ext != ".jpeg" {
			p.Filename = strings.TrimSuffix(p.Filename, ext) + ".jpg"
		}
	case "image/webp":
		if ext != ".webp" {
			p.Filename = strings.TrimSuffix(p.Filename, ext) + ".webp"
		}
	}
	if p.Size == 0 {
		p.Size = int64(len(p.Data))
	}
	return nil
}

// Create a new instrument photo
func AddInstrumentPhoto(db *gorm.DB, photo *InstrumentPhoto) error { return db.Create(photo).Error }

// Get a single instrument photo by id
func GetInstrumentPhoto(db *gorm.DB, id uint) (*InstrumentPhoto, error) {
	var p InstrumentPhoto
	if err := db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// List photos for an instrument
func ListInstrumentPhotos(db *gorm.DB, instrumentID uint) ([]InstrumentPhoto, error) {
	var out []InstrumentPhoto
	if err := db.Where("instrument_id = ?", instrumentID).Order("id desc").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

// Delete an instrument photo
func DeleteInstrumentPhoto(db *gorm.DB, id uint) error {
	return db.Delete(&InstrumentPhoto{}, id).Error
}
