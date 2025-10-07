package models

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

// StationPhoto stores a station's photo binary in the database
// Allowed formats: PNG, JPEG, WebP
// TableName: StationPhotos

type StationPhoto struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	STNID       uint      `gorm:"not null;index" json:"stn_id"`
	Filename    string    `gorm:"size:255;not null" json:"filename"`
	ContentType string    `gorm:"size:100;not null" json:"content_type"`
	Size        int64     `gorm:"not null" json:"size"`
	Data        []byte    `gorm:"type:longblob;not null" json:"-"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StationPhoto) TableName() string { return "StationPhotos" }

var allowedPhotoMIMEs = map[string]struct{}{
	"image/png":  {},
	"image/jpeg": {},
	"image/webp": {},
}

func (p *StationPhoto) BeforeCreate(tx *gorm.DB) error { return p.validate() }
func (p *StationPhoto) BeforeUpdate(tx *gorm.DB) error { return p.validate() }

func (p *StationPhoto) validate() error {
	if p.STNID == 0 {
		return errors.New("stn id is required")
	}
	if len(p.Data) == 0 {
		return errors.New("image data is required")
	}
	if p.ContentType == "" {
		p.ContentType = http.DetectContentType(p.Data)
	}
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

func ensureExt(name, ext string) string {
	cur := strings.ToLower(filepath.Ext(name))
	if cur == ext {
		return name
	}
	return strings.TrimSuffix(name, cur) + ext
}

// Create a new station photo
func AddStationPhoto(db *gorm.DB, photo *StationPhoto) error { return db.Create(photo).Error }

// Get a single photo by id
func GetStationPhoto(db *gorm.DB, id uint) (*StationPhoto, error) {
	var p StationPhoto
	if err := db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// List photos for a station
func ListStationPhotosBySTN(db *gorm.DB, stnID uint) ([]StationPhoto, error) {
	var out []StationPhoto
	if err := db.Where("stn_id = ?", stnID).Order("id desc").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

// Delete a photo
func DeleteStationPhoto(db *gorm.DB, id uint) error { return db.Delete(&StationPhoto{}, id).Error }
