package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshTokenModel struct {
	UUID                   uuid.UUID  `db:"uuid"`
	UserUUID               uuid.UUID  `db:"user_uuid"`
	TokenHash              []byte     `db:"token_hash"`
	AccessTokenHash        []byte     `db:"access_token_hash"`
	ExpiredDate            *time.Time `db:"expired_date"`
	AccessTokenExpiredDate *time.Time `db:"access_token_expired_date"`
	CreatedDate            time.Time  `db:"created_date"`
	RevokeDate             *time.Time `db:"revoke_date"`
	ReplacedBy             uuid.UUID  `db:"replaced_by"`
	UserAgent              string     `db:"user_agent"`
	IpAddress              string     `db:"ip_address"`
}
