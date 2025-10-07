package models

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserAvatar stores a user's avatar binary separately from Users table
// Allowed formats: PNG, JPEG, WebP
// TableName: UserAvatars

type UserAvatar struct {
	UserID      uint      `gorm:"primaryKey;index" json:"user_id"`
	Data        []byte    `gorm:"type:longblob" json:"-"`
	ContentType string    `gorm:"size:100" json:"content_type"`
	Size        int64     `json:"size"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (UserAvatar) TableName() string { return "UserAvatars" }

var allowedAvatarMIMEs = map[string]struct{}{
	"image/png":  {},
	"image/jpeg": {},
	"image/webp": {},
}

// UpsertUserAvatar validates and stores the avatar binary for a user
func UpsertUserAvatar(db *gorm.DB, userID uint, data []byte, contentType string) error {
	if len(data) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "avatar image is required")
	}
	ct := contentType
	if ct == "" {
		ct = http.DetectContentType(data)
	}
	if _, ok := allowedAvatarMIMEs[ct]; !ok {
		return fiber.NewError(fiber.StatusBadRequest, "unsupported avatar type; allowed: PNG, JPEG, WebP")
	}
	ua := UserAvatar{UserID: userID, Data: data, ContentType: ct, Size: int64(len(data))}
	return db.Save(&ua).Error
}

// RemoveUserAvatar deletes the avatar row for a user
func RemoveUserAvatar(db *gorm.DB, userID uint) error {
	return db.Delete(&UserAvatar{}, "user_id = ?", userID).Error
}

// GetUserAvatar retrieves a user's avatar
func GetUserAvatar(db *gorm.DB, userID uint) (*UserAvatar, error) {
	var ua UserAvatar
	if err := db.First(&ua, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &ua, nil
}
