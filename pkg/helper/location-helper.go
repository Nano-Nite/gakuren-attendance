package helper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gakuren-system.com/pkg/db"
	"gakuren-system.com/pkg/model/attendance"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func ValidateRequestLocation(data attendance.CreateLocationRequest) error {
	query := `
	select 
		code
		,name
	from attendance_sch.attendance_location
	where lower(code) = lower($1)
	OR lower(name) = lower($2)
	`

	result, err := db.GetMultipleDataByQuery[attendance.LocationModel](query,
		data.Code,
		data.Name,
	)
	if err != nil {
		return err
	}
	if len(*result) > 0 {
		return errors.New("Multiple location data found")
	}

	return nil
}

func CreateLocation(data attendance.LocationModel, tx pgx.Tx, ctx context.Context) error {
	query := `
		insert into attendance_sch.attendance_location
			(tenant_uuid, school_uuid, code, name, latitude, longitude, radius_meter, is_active, created_date, created_by, updated_date, updated_by )
		values
			($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		returning uuid;
	`
	var id uuid.UUID

	err := tx.QueryRow(ctx, query,
		data.TenantUUID,
		data.SchoolUUID,
		data.Code,
		data.Name,
		data.Latitude,
		data.Longitude,
		data.RadiusMeter,
		true,
		time.Now(),
		data.CreatedBy,
		nil,
		nil,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("insert location: %w", err)
	}

	return nil
}

func GetAllLocation(schoolUUID, tenantUUID uuid.UUID) ([]attendance.GetAllLocationModel, error) {
	query := `
		select
			uuid
			,code
			,name
			,latitude
			,longitude
			,radius_meter
		from attendance_sch.attendance_location
		where tenant_uuid = $1
		and school_uuid = $2
	`
	rows, err := db.GetMultipleDataByQuery[attendance.GetAllLocationModel](query, tenantUUID, schoolUUID)
	if err != nil {
		return nil, err
	}

	return *rows, nil
}
