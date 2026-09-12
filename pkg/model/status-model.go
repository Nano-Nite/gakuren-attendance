package model

import (
	"time"

	"github.com/google/uuid"
)

type StatusModel struct {
	UUID        uuid.UUID              `db:"uuid"`
	Name        string                 `db:"name"`
	AbbrName    string                 `db:"abbr_name"`
	Code        string                 `db:"code"`
	Category    string                 `db:"category"`
	Description string                 `db:"description"`
	SortOrder   string                 `db:"sort_order"`
	IsActive    bool                   `db:"is_active"`
	IsTerminal  bool                   `db:"is_terminal"`
	Metadata    map[string]interface{} `db:"metadata"`
	CreatedDate time.Time              `db:"created_date"`
	UpdatedDate *time.Time             `db:"updated_date"`
}
