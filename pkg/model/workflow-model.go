package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ApprovalWorkflow struct {
	UUID        uuid.UUID  `db:"uuid"`
	TenantUUID  uuid.UUID  `db:"tenant_uuid"`
	MenuUUID    uuid.UUID  `db:"menu_uuid"`
	SchoolUUID  uuid.UUID  `db:"school_uuid"`
	Code        string     `db:"code"`
	Name        string     `db:"name"`
	ActionCode  string     `db:"action_code"`
	StatusUUID  uuid.UUID  `db:"status_uuid"`
	CreatedBy   uuid.UUID  `db:"created_by"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedBy   *uuid.UUID `db:"updated_by"`
	UpdatedDate *time.Time `db:"updated_date"`
}

type ApprovalInstance struct {
	UUID                 uuid.UUID       `db:"uuid" json:"uuid"`
	ApprovalWorkflowUUID uuid.UUID       `db:"approval_workflow_uuid" json:"approval_workflow_uuid"`
	TenantUUID           uuid.UUID       `db:"tenant_uuid" json:"tenant_uuid"`
	SchoolUUID           uuid.UUID       `db:"school_uuid" json:"school_uuid"`
	TicketNumber         string          `db:"ticket_number" json:"ticket_number"`
	EntityType           string          `db:"entity_type" json:"entity_type"`
	EntityUUID           *uuid.UUID      `db:"entity_uuid" json:"entity_uuid,omitempty"`
	ActionCode           string          `db:"action_code" json:"action_code"`
	RequestData          json.RawMessage `db:"request_data" json:"request_data"`
	CurrentStep          int             `db:"current_step" json:"current_step"`
	StatusUUID           uuid.UUID       `db:"status_uuid" json:"status_uuid"`
	RequestedBy          uuid.UUID       `db:"requested_by" json:"requested_by"`
	RequestedDate        time.Time       `db:"requested_date" json:"requested_date"`
	FinalizedBy          *uuid.UUID      `db:"finalized_by" json:"finalized_by,omitempty"`
	FinalizedDate        *time.Time      `db:"finalized_date" json:"finalized_date,omitempty"`
	UpdatedDate          *time.Time      `db:"updated_date" json:"updated_date,omitempty"`
}

type ApprovalAction struct {
	UUID                 uuid.UUID  `db:"uuid" json:"uuid"`
	ApprovalInstanceUUID uuid.UUID  `db:"approval_instance_uuid" json:"approval_instance_uuid"`
	ApprovalStepUUID     *uuid.UUID `db:"approval_step_uuid" json:"approval_step_uuid,omitempty"`
	ActionCode           string     `db:"action_code" json:"action_code"`
	ActedBy              uuid.UUID  `db:"acted_by" json:"acted_by"`
	Note                 *string    `db:"note" json:"note,omitempty"`
	CreatedDate          time.Time  `db:"created_date" json:"created_date"`
}

type ApprovalStep struct {
	UUID                 uuid.UUID  `db:"uuid" json:"uuid"`
	ApprovalWorkflowUUID uuid.UUID  `db:"approval_workflow_uuid" json:"approval_workflow_uuid"`
	StepOrder            int        `db:"step_order" json:"step_order"`
	ApproverRoleUUID     uuid.UUID  `db:"approver_role_uuid" json:"approver_role_uuid"`
	PermissionUUID       *uuid.UUID `db:"permission_uuid" json:"permission_uuid,omitempty"`
	RequiredApprovals    int        `db:"required_approvals" json:"required_approvals"`
	CreatedBy            uuid.UUID  `db:"created_by" json:"created_by"`
	CreatedDate          time.Time  `db:"created_date" json:"created_date"`
	UpdatedBy            *uuid.UUID `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedDate          *time.Time `db:"updated_date" json:"updated_date,omitempty"`
}
