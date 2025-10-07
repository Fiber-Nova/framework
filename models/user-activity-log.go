package models

import (
	"time"

	"gorm.io/gorm"
)

// UserActivityLog represents a user activity record
type UserActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Activity  string    `gorm:"size:255;not null" json:"activity"`
	Details   string    `gorm:"type:text" json:"details"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	UserAgent string    `gorm:"size:512" json:"user_agent"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName sets the table name
func (UserActivityLog) TableName() string {
	return "UserActivityLogs"
}

// LogActivity creates a new activity log
func LogActivity(db *gorm.DB, userID uint, activity, details, ip, userAgent string) error {
	log := &UserActivityLog{
		UserID:    userID,
		Activity:  activity,
		Details:   details,
		IPAddress: ip,
		UserAgent: userAgent,
	}
	return db.Create(log).Error
}

// GetUserActivities fetches activities for a user with pagination
func GetUserActivities(db *gorm.DB, userID uint, page, pageSize int) ([]UserActivityLog, int64, error) {
	var logs []UserActivityLog
	var count int64

	offset := (page - 1) * pageSize
	err := db.Model(&UserActivityLog{}).
		Where("user_id = ?", userID).
		Count(&count).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error
	return logs, count, err
}

// GetActivitiesByTimeRange fetches activities within a time range
func GetActivitiesByTimeRange(db *gorm.DB, start, end time.Time) ([]UserActivityLog, error) {
	var logs []UserActivityLog
	err := db.Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetActivityStats gets statistics by activity type
func GetActivityStats(db *gorm.DB, days int) (map[string]int64, error) {
	stats := make(map[string]int64)
	var results []struct {
		Activity string
		Count    int64
	}

	startTime := time.Now().AddDate(0, 0, -days)
	err := db.Model(&UserActivityLog{}).
		Select("activity, COUNT(*) as count").
		Where("created_at >= ?", startTime).
		Group("activity").
		Scan(&results).Error

	for _, r := range results {
		stats[r.Activity] = r.Count
	}
	return stats, err
}

// CleanupOldLogs removes logs older than specified days
func CleanupOldLogs(db *gorm.DB, days int) (int64, error) {
	result := db.Where("created_at < ?", time.Now().AddDate(0, 0, -days)).
		Delete(&UserActivityLog{})
	return result.RowsAffected, result.Error
}
