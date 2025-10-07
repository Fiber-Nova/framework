package models

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

var Store = session.New()

// User model
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Email     string    `gorm:"size:255" json:"email"`
	Active    bool      `gorm:"default:true" json:"active"`
	Deleted   string    `gorm:"type:ENUM('Yes','No');default:'No'" json:"deleted"`
	Role      string    `gorm:"type:ENUM('Root','Admin','Inspector');default:'Inspector';index" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName sets the table name to 'Users' for GORM
func (User) TableName() string {
	return "Users"
}

// CreateUser creates a new user in the database with validation and transaction
func CreateUser(db *gorm.DB, user *User) error {
	// Validate required fields
	if user.Username == "" || user.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "username and password are required")
	}
	// Validate email format if provided
	if user.Email != "" && !strings.Contains(user.Email, "@") {
		return fiber.NewError(fiber.StatusBadRequest, "invalid email format")
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if username already exists
	var count int64
	if err := tx.Model(&User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "failed to check username availability")
	}
	if count > 0 {
		tx.Rollback()
		return fiber.NewError(fiber.StatusConflict, "username already exists")
	}

	// Set created_at timestamp
	user.CreatedAt = time.Now()

	// Create user
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create user")
	}

	return tx.Commit().Error
}

// DeactivateUser sets the user's Active field to false
func DeactivateUser(db *gorm.DB, userID uint) error {
	return db.Model(&User{}).Where("id = ?", userID).Update("active", false).Error
}

// MarkUserDeleted sets the user's Deleted flag to "Yes"
func MarkUserDeleted(db *gorm.DB, userID uint) error {
	return db.Model(&User{}).Where("id = ?", userID).Update("deleted", "Yes").Error
}

// DeleteUser deletes a user by ID
func DeleteUser(db *gorm.DB, userID uint) error {
	return db.Delete(&User{}, userID).Error
}

// FindUserByID retrieves a user by ID
func FindUserByID(db *gorm.DB, userID uint) (*User, error) {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetLoggedInUser retrieves the logged-in user model from the session
func GetLoggedInUser(c *fiber.Ctx, db *gorm.DB) (*User, error) {
	sess, err := Store.Get(c)
	if err != nil {
		return nil, err
	}
	userIDVal := sess.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}
	user, err := FindUserByID(db, uint(userID))
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Note: Avatar-related fields and helpers were moved to models/avatar.go (UserAvatar).
