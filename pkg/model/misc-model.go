package model

import (
	"time"

	"github.com/google/uuid"
)

type CountResult struct {
	Count int `db:"count"`
}

type DataStatistics struct {
	CurrentRow  int `json:"current_row"`
	StartRow    int `json:"start_row"`
	EndRow      int `json:"end_row"`
	TotalRow    int `json:"total_row"`
	CurrentPage int `json:"current_page"`
	MaxPage     int `json:"max_page"`
	RowPerPage  int `json:"row_per_page"`
}

type SearchPayload struct {
	Search     *string                   `json:"search,omitempty"`
	Filter     *map[string]interface{}   `json:"filter,omitempty"`
	Page       *int                      `json:"page,omitempty"`
	RowPerPage *int                      `json:"row_per_page,omitempty"`
	SortBy     *[]map[string]interface{} `json:"sort_by,omitempty"`
}

type GetPermission struct {
	Code string `db:"code"`
}

type GetMenuUUID struct {
	UUID uuid.UUID `db:"uuid"`
}

type VariableModel struct {
	UUID        uuid.UUID  `db:"uuid"`
	Key         string     `db:"key"`
	Value       string     `db:"value"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate *time.Time `db:"updated_date"`
}

type ExecuteApprovalPayload struct {
	UUID    string `db:"uuid" json:"uuid"`
	Command string `db:"command" json:"command"`
	Note    string `db:"note" json:"note"`
}

type ReadGenderResult struct {
	UUID     uuid.UUID `db:"uuid" json:"uuid"`
	Name     string    `db:"name" json:"name"`
	AbbrName *string   `db:"abbr_name" json:"abbr_name"`
	Status   string    `db:"status" json:"status"`
}

type ResultPositionModel struct {
	UUID     uuid.UUID `db:"uuid" json:"uuid"`
	Name     string    `db:"name" json:"name"`
	AbbrName *string   `db:"abbr_name" json:"abbr_name"`
	Status   string    `db:"status" json:"status"`
}

type ResultTitleModel struct {
	UUID     uuid.UUID `db:"uuid" json:"uuid"`
	Name     string    `db:"name" json:"name"`
	AbbrName *string   `db:"abbr_name" json:"abbr_name"`
	Status   string    `db:"status" json:"status"`
	IsPrefix bool      `db:"is_prefix" json:"is_prefix"`
	Sequence int       `db:"sequence" json:"sequence"`
}

type ResultSubjectModel struct {
	UUID     uuid.UUID `db:"uuid" json:"uuid"`
	Name     string    `db:"name" json:"name"`
	AbbrName *string   `db:"abbr_name" json:"abbr_name"`
}
