package models

import "time"

type AuditLog struct {
	ID         string
	TenantID   string
	UserID     *string
	Action     string
	Resource   string
	ResourceID *string
	Method     string
	Path       string
	StatusCode int
	IPAddress  *string
	UserAgent  string
	Metadata   []byte
	CreatedAt  time.Time
}
