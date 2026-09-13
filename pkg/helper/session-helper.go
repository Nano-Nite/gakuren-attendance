package helper

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"gakuren-system.com/pkg/db"
	"gakuren-system.com/pkg/model/attendance"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func ValidateRequestSession(data attendance.CreateSession, schoolUUID, tenantUUID uuid.UUID) error {
	query := `
	select
		*
	from attendance_sch.attendance_session t 
	where t.tenant_uuid = $1
		and t.school_uuid = $2
		and t.location_uuid = $3
		and ( t.status  = $4 or t.closed_by is not null or t.closed_date is not null)
	`

	result, err := db.GetMultipleDataByQuery[attendance.SessionModel](query,
		tenantUUID,
		schoolUUID,
		data.LocatoinUUID,
		attendance.SESSION_STATUS_ACTIVE,
	)
	if err != nil {
		return err
	}
	if len(*result) > 0 {
		return errors.New("Multiple location data found")
	}

	return nil
}

func CreateSession(data attendance.SessionModel, tx pgx.Tx, ctx context.Context) error {
	query := `
		insert into attendance_sch.attendance_session
			(tenant_uuid, school_uuid, session_code, attendance_type, target_type, location_uuid, valid_from, valid_until, status, created_date, created_by)
		values
			($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		returning uuid;
	`
	var id uuid.UUID

	err := tx.QueryRow(ctx, query,
		data.TenantUUID,
		data.SchoolUUID,
		GenerateSessionCode(),
		data.AttendanceType,
		data.TargetType,
		data.LocatoinUUID,
		data.ValidFrom,
		data.ValidUntil,
		data.Status,
		time.Now(),
		data.CreatedBy,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("insert location: %w", err)
	}

	return nil
}

func GetActiveSession(schoolUUID, tenantUUID uuid.UUID) (*attendance.GetActiveSessionModel, error) {
	query := `
		select
			uuid session_uuid
			,session_code
			,status
			,'' as qr_token
			,attendance_type
			,valid_from
			,valid_until
			,jsonb_build_object(
				'name', (select al."name"  from attendance_sch.attendance_location al where al.uuid = s.location_uuid),
				'geofence_radius_meter', (select al.radius_meter from attendance_sch.attendance_location al where al.uuid = s.location_uuid)
			) location
			, jsonb_build_object(
				'name', (select u."name"  from user_sch."user" u where u."uuid" = s.created_by)
			) as created_by
		from attendance_sch.attendance_session s
		where tenant_uuid = $1
		and school_uuid = $2;
	`
	result, err := db.GetSingleDataByQuery[attendance.GetActiveSessionModel](query, tenantUUID, schoolUUID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func generateRandomChar(n int) string {
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(attendance.CHARSET_SESSION_CODE))))
		b[i] = attendance.CHARSET_SESSION_CODE[num.Int64()]
	}
	return string(b)
}

func GenerateSessionCode() string {
	prefix := "ABS"
	date := time.Now().Format("20060102") // Format YYYYMMDD di Go
	randomStr := generateRandomChar(6)    // 6 karakter/digit acak

	return fmt.Sprintf("%s-%s-%s", prefix, date, randomStr)
}
