package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Unified audit log for all entities
// TableName: AuditLogs
// Entity denotes which table/model was modified (e.g., "user","stn","instrument","store","instrument_type")
// Action: e.g., "create","update","delete","login"
// Changes: optional JSON payload with diff/fields

type AuditLog struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Entity    string          `gorm:"size:64;index:idx_entity,priority:1" json:"entity"`
	EntityID  uint            `gorm:"index:idx_entity,priority:2" json:"entity_id"`
	Action    string          `gorm:"size:32;index" json:"action"`
	ActorID   *uint           `gorm:"index" json:"actor_id,omitempty"`
	Actor     string          `gorm:"size:100" json:"actor,omitempty"`
	Details   string          `json:"details,omitempty"`
	Changes   json.RawMessage `json:"changes,omitempty"`
	CreatedAt time.Time       `gorm:"autoCreateTime;index:idx_entity,priority:3" json:"created_at"`
}

func (AuditLog) TableName() string { return "AuditLogs" }

// LogAudit inserts a new audit entry
func LogAudit(db *gorm.DB, entity string, entityID uint, action string, actorID *uint, actor string, details string, changes any) error {
	var payload []byte
	if changes != nil {
		if b, err := json.Marshal(changes); err == nil {
			payload = b
		}
	}
	entry := AuditLog{
		Entity: entity, EntityID: entityID, Action: action,
		ActorID: actorID, Actor: actor, Details: details, Changes: payload,
	}
	return db.Create(&entry).Error
}
