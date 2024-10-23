package audit_log

import "time"

type AuditLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
