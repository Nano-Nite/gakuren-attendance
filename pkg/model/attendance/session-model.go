package attendance

import (
	"time"

	"github.com/google/uuid"
)

// SESSION STATUS VAR
const SESSION_STATUS_ACTIVE = "ACTIVE"
const SESSION_STATUS_CLOSED = "CLOSED"
const SESSION_STATUS_EXPIRED = "EXPIRED"
const SESSION_STATUS_CANCELLED = "CANCELLED"

const CHARSET_SESSION_CODE = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type SessionModel struct {
	TenantUUID     uuid.UUID  `json:"tenant_uuid"`
	SchoolUUID     uuid.UUID  `json:"school_uuid"`
	SessionCode    string     `json:"session_code"`
	AttendanceType string     `json:"attendance_type"`
	TargetType     string     `json:"target_type"`
	LocatoinUUID   uuid.UUID  `json:"location_uuid"`
	ValidFrom      time.Time  `json:"valid_from"`
	ValidUntil     time.Time  `json:"valid_until"`
	Status         string     `json:"status"`
	QR             string     `json:"qr"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	CreatedDate    time.Time  `json:"created_date"`
	ClosedBy       *uuid.UUID `json:"closed_by"`
	ClosedDate     *time.Time `json:"closed_date"`
	UpdatedBy      *uuid.UUID `json:"updated_by"`
	UpdatedDate    *time.Time `json:"updated_date"`
}

type CreateSession struct {
	LocatoinUUID   uuid.UUID `json:"location_uuid"`
	AttendanceType string    `json:"attendance_type"`
	TargetType     string    `json:"target_type"`
	ValidFrom      time.Time `json:"valid_from"`
	ValidUntil     time.Time `json:"valid_until"`
}

type GetActiveSessionModel struct {
	SessionUUID    uuid.UUID   `db:"session_uuid" json:"session_uuid"`
	SessionCode    string      `db:"session_code" json:"session_code"`
	Status         string      `db:"status" json:"status"`
	QRToken        *string     `db:"qr_token" json:"qr_token"`
	AttendanceType string      `db:"attendance_type" json:"attendance_type"`
	ValidFrom      time.Time   `db:"valid_from" json:"valid_from"`
	ValidUntil     time.Time   `db:"valid_until" json:"valid_until"`
	Location       interface{} `db:"location" json:"location"`
	CreatedBy      interface{} `db:"created_by" json:"created_by"`
}

type QRPayload struct {
	Type        string    `db:"type" json:"type"`
	SessionCode string    `db:"session_code" json:"session_code"`
	TenantUUID  uuid.UUID `db:"tenant_uuid" json:"tenant_uuid"`
	SchoolUUID  uuid.UUID `db:"school_uuid" json:"school_uuid"`
}

type SignedQR struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}
