package domain

import (
	"github.com/google/uuid"
	"time"
)

type AuditLogEntry struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Action    string    `json:"action"`
	Token     string    `json:"token"`
	Kind      *Kind     `json:"kind"`
	CreatedAt time.Time `json:"created_at"`
}
