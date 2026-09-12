package attendance

import (
	"time"

	"github.com/google/uuid"
)

type LocationModel struct {
	UUID        uuid.UUID  `json:"uuid"`
	TenantUUID  uuid.UUID  `json:"tenant_uuid"`
	SchoolUUID  uuid.UUID  `json:"school_uuid"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Latitude    *float64   `json:"latitude,omitempty"`
	Longitude   *float64   `json:"longitude,omitempty"`
	RadiusMeter int        `json:"radius_meter"`
	IsActive    bool       `json:"is_active"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	CreatedDate time.Time  `json:"created_date"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
	UpdatedDate *time.Time `json:"updated_date,omitempty"`
}

type GetAllLocationModel struct {
	UUID        uuid.UUID `json:"uuid"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	RadiusMeter int       `json:"radius_meter"`
}

type CreateLocationRequest struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	RadiusMeter int      `json:"radius_meter"`
}

type UpdateLocationRequest struct {
	Code        *string  `json:"code"`
	Name        *string  `json:"name"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	RadiusMeter *int     `json:"radius_meter"`
	IsActive    *bool    `json:"is_active"`
}
