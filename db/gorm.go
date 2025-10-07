package db

import (
	"backend-meta-data/config"
	"backend-meta-data/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectGormDB(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// InitDBIfNeeded centralizes schema migrations and one-off adjustments
func InitDBIfNeeded(db *gorm.DB) error {
	// Allow disabling auto-migration via environment flag
	if os.Getenv("AUTO_MIGRATE") == "false" {
		return nil
	}

	// Auto-migrate all core models here
	if err := db.AutoMigrate(
		&models.Configuration{},
		&models.User{},
		&models.StationType{},
		&models.Station{},
		&models.AuditLog{},
		&models.InstrumentType{},
		&models.Instrument{},
		&models.InspectionRecord{},
		&models.InspectionFormAttachment{},
		&models.MaintenanceNotice{},
		&models.InventoryStore{},
		&models.StationPhoto{},
		&models.InstrumentPhoto{},
		&models.UserAvatar{},
		&models.MaintNoticeTemplate{},
	); err != nil {
		return err
	}

	// Seed maintenance notice templates from filesystem if table is empty
	var count int64
	if err := db.Model(&models.MaintNoticeTemplate{}).Count(&count).Error; err == nil && count == 0 {
		dirs := []string{"./src/backend/email_templates/maint_notice", "./backend/email_templates/maint_notice", "./email_templates/maint_notice"}
		for _, d := range dirs {
			_ = models.SeedMaintNoticeTemplatesFromFS(db, d)
		}
	}

	// Ensure StationTypeID column exists on Stations and create FK constraint (only if missing)
	if db.Migrator().HasTable(&models.Station{}) {
		if !db.Migrator().HasColumn(&models.Station{}, "StationTypeID") {
			if err := db.Migrator().AddColumn(&models.Station{}, "StationTypeID"); err != nil {
				return err
			}
		}
		if !db.Migrator().HasConstraint(&models.Station{}, "StationType") {
			if err := db.Migrator().CreateConstraint(&models.Station{}, "StationType"); err != nil {
				return err
			}
		}
	}

	// InstrumentTypes: ensure Code values exist and are unique before unique index is enforced
	if db.Migrator().HasTable(&models.InstrumentType{}) {
		if !db.Migrator().HasColumn(&models.InstrumentType{}, "Code") {
			if err := db.Migrator().AddColumn(&models.InstrumentType{}, "Code"); err != nil {
				return err
			}
		}
		// Backfill empty or NULL codes to unique values (e.g., IT{id}) to satisfy unique constraint
		if err := db.Exec("UPDATE `InstrumentTypes` SET `code` = CONCAT('IT', id) WHERE `code` IS NULL OR `code` = ''").Error; err != nil {
			return fmt.Errorf("backfill instrument type codes failed: %w", err)
		}
	}

	// Ensure Users.role ENUM matches allowed roles
	if db.Migrator().HasTable(&models.User{}) {
		// Try altering column to ENUM if DB supports it (MySQL)
		_ = db.Exec("ALTER TABLE `Users` MODIFY COLUMN `role` ENUM('Root','Admin','Inspector') NOT NULL DEFAULT 'Inspector'").Error
	}
	return nil
}
