package helper

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
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
		and t.status = $4 and t.closed_by is null and t.closed_date is null
		and t.valid_until > now()
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
		return errors.New("Multiple session data found")
	}

	return nil
}

func CreateSession(data *attendance.SessionModel, tx pgx.Tx, ctx context.Context) error {
	sessionCode := GenerateSessionCode()

	qr, err := GenerateSessionQR(attendance.QRPayload{
		Type:        "attendance_session",
		SessionCode: sessionCode,
		TenantUUID:  data.TenantUUID,
		SchoolUUID:  data.SchoolUUID,
	})

	query := `
		insert into attendance_sch.attendance_session
			(tenant_uuid, school_uuid, session_code, attendance_type, target_type, location_uuid, valid_from, valid_until, status, qr, created_date, created_by)
		values
			($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		returning uuid;
	`
	var id uuid.UUID

	err = tx.QueryRow(ctx, query,
		data.TenantUUID,
		data.SchoolUUID,
		sessionCode,
		data.AttendanceType,
		data.TargetType,
		data.LocatoinUUID,
		data.ValidFrom,
		data.ValidUntil,
		data.Status,
		qr,
		time.Now(),
		data.CreatedBy,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}

	return err
}

// GenerateSessionQR signs the payload; Base64 encoding does not hide its contents.
func GenerateSessionQR(payload attendance.QRPayload) (string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("create qr session: %w", err)
	}

	private, err := base64.StdEncoding.DecodeString(strings.TrimSpace(os.Getenv("ED25519_PRIVATE_KEY")))
	if err != nil {
		return "", fmt.Errorf("decode ED25519_PRIVATE_KEY: %w", err)
	}
	switch len(private) {
	case ed25519.SeedSize:
		private = ed25519.NewKeyFromSeed(private)
	case ed25519.PrivateKeySize:
		if !bytes.Equal(private, ed25519.NewKeyFromSeed(private[:ed25519.SeedSize])) {
			return "", errors.New("ED25519_PRIVATE_KEY has an inconsistent public key")
		}
	default:
		return "", errors.New("ED25519_PRIVATE_KEY must be standard Base64 of a 32-byte seed or 64-byte private key")
	}

	signature := ed25519.Sign(
		private,
		payloadBytes,
	)

	qr := attendance.SignedQR{
		Payload:   base64.RawURLEncoding.EncodeToString(payloadBytes),
		Signature: base64.RawURLEncoding.EncodeToString(signature),
	}

	qrBytes, err := json.Marshal(qr)
	if err != nil {
		return "", fmt.Errorf("encode signed qr: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(qrBytes), nil
}

// VerifySessionQR verifies authenticity only. The caller must also check the
// session's database status, validity window and tenant/school authorization.
// publicKey must come from trusted configuration, never from the scanned QR.
func VerifySessionQR(token string, publicKey ed25519.PublicKey) (*attendance.QRPayload, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return nil, errors.New("invalid Ed25519 public key size")
	}
	qrBytes, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("decode qr token: %w", err)
	}
	var qr attendance.SignedQR
	if err := json.Unmarshal(qrBytes, &qr); err != nil {
		return nil, fmt.Errorf("decode signed qr: %w", err)
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(qr.Payload)
	if err != nil {
		return nil, fmt.Errorf("decode qr payload: %w", err)
	}
	signature, err := base64.RawURLEncoding.DecodeString(qr.Signature)
	if err != nil {
		return nil, fmt.Errorf("decode qr signature: %w", err)
	}
	if !ed25519.Verify(publicKey, payloadBytes, signature) {
		return nil, errors.New("invalid qr signature")
	}
	var payload attendance.QRPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("decode qr payload: %w", err)
	}
	if payload.Type != "attendance_session" || payload.SessionCode == "" || payload.TenantUUID == uuid.Nil || payload.SchoolUUID == uuid.Nil {
		return nil, errors.New("invalid attendance qr payload")
	}
	return &payload, nil
}

func GetActiveSession(schoolUUID, tenantUUID uuid.UUID) (*attendance.GetActiveSessionModel, error) {
	query := `
		select
			uuid session_uuid
			,session_code
			,status
			,qr as qr_token
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
		and school_uuid = $2
		and status = 'ACTIVE' and closed_by is null and closed_date is null
		and valid_from <= now() and valid_until > now()
		order by created_date desc
		limit 1;
	`
	result, err := db.GetSingleDataByQuery[attendance.GetActiveSessionModel](query, tenantUUID, schoolUUID)
	if err != nil {
		return nil, err
	}

	token, err := GenerateSessionQR(attendance.QRPayload{
		Type: "attendance_session", SessionCode: result.SessionCode,
		TenantUUID: tenantUUID, SchoolUUID: schoolUUID,
	})
	if err != nil {
		return nil, err
	}
	result.QRToken = &token
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
