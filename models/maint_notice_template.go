package models

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MaintNoticeTemplate struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:150;uniqueIndex;not null" json:"name"` // e.g. default.html
	Label       string    `gorm:"size:200;not null" json:"label"`
	Description string    `json:"description"`
	Engine      string    `gorm:"type:ENUM('html','quill','gohtml');default:'html'" json:"engine"`
	Content     string    `gorm:"type:longtext" json:"content"`
	Active      bool      `gorm:"default:true;index" json:"active"`
	Deleted     string    `gorm:"type:ENUM('Yes','No');default:'No';index" json:"deleted"`
	Version     int       `gorm:"default:1" json:"version"`
	CreatedByID *uint     `gorm:"index" json:"created_by_id,omitempty"`
	UpdatedByID *uint     `gorm:"index" json:"updated_by_id,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (MaintNoticeTemplate) TableName() string { return "MaintNoticeTemplates" }

func UpsertMaintNoticeTemplate(db *gorm.DB, t *MaintNoticeTemplate) error {
	if strings.TrimSpace(t.Name) == "" {
		return gorm.ErrInvalidData
	}
	var existing MaintNoticeTemplate
	if err := db.Where("name = ?", t.Name).First(&existing).Error; err == nil {
		if strings.TrimSpace(t.Label) != "" {
			existing.Label = t.Label
		}
		if t.Description != "" {
			existing.Description = t.Description
		}
		if strings.TrimSpace(t.Engine) != "" {
			existing.Engine = t.Engine
		}
		if t.Content != "" {
			existing.Content = t.Content
		}
		existing.Active = t.Active || existing.Active
		existing.Deleted = "No"
		existing.Version = existing.Version + 1
		return db.Save(&existing).Error
	}
	if strings.TrimSpace(t.Label) == "" {
		t.Label = t.Name
	}
	if strings.TrimSpace(t.Engine) == "" {
		t.Engine = "html"
	}
	if strings.TrimSpace(t.Deleted) == "" {
		t.Deleted = "No"
	}
	return db.Create(t).Error
}

func GetMaintNoticeTemplateByName(db *gorm.DB, name string) (*MaintNoticeTemplate, error) {
	var t MaintNoticeTemplate
	if err := db.Where("name = ? AND deleted = 'No' AND active = 1", name).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func ListMaintNoticeTemplates(db *gorm.DB) ([]MaintNoticeTemplate, error) {
	var list []MaintNoticeTemplate
	if err := db.Where("deleted = 'No' AND active = 1").Order("label ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// SeedMaintNoticeTemplatesFromFS loads templates from a folder into DB if missing
func SeedMaintNoticeTemplatesFromFS(db *gorm.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	reName := regexp.MustCompile(`(?is)Template\s*Name:\s*(.+?)\s*(?:\r?\n|-->|\r)`) // extract name from HTML comment
	reDesc := regexp.MustCompile(`(?is)Description:\s*(.+?)\s*(?:\r?\n|-->|\r)`)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".html") && !strings.HasSuffix(strings.ToLower(name), ".htm") {
			continue
		}
		path := filepath.Join(dir, name)
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(b)
		label := name
		desc := ""
		if m := reName.FindStringSubmatch(content); len(m) == 2 {
			label = strings.TrimSpace(m[1])
		}
		if m := reDesc.FindStringSubmatch(content); len(m) == 2 {
			desc = strings.TrimSpace(m[1])
		}
		_ = UpsertMaintNoticeTemplate(db, &MaintNoticeTemplate{Name: name, Label: label, Description: desc, Engine: "html", Content: content, Active: true})
	}
	return nil
}
