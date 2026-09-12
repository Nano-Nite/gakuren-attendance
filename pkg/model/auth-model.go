package model

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginPayload struct {
	Email string `json:"email,omitempty"`
}

type LogoutPayload struct {
	Email string `json:"email,omitempty"`
}

type PublicKeyResponse struct {
	PublicKey string `json:"public_key,omitempty"`
}

type JWTService struct {
	PrivateKey *rsa.PrivateKey
}

type AccessTokenClaims struct {
	SessionID  string   `json:"sid"`
	TenantUUID string   `json:"tenant_uuid"`
	SchoolUUID string   `json:"school_uuid"`
	Permission []string `json:"permission"`
	Username   string   `json:"username"`
	Roles      []string `json:"roles,omitempty"`

	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Username string `json:"username"`
	TokenID  string `json:"token_id"`
	jwt.RegisteredClaims
}

type LoginData struct {
	UUID          uuid.UUID `json:"uuid"`
	TenantUUID    uuid.UUID `json:"tenant_uuid"`
	SchoolUUID    uuid.UUID `json:"school_uuid"`
	TenantName    string    `json:"tenant_name"`
	Code          string    `json:"code"`
	Timezone      string    `json:"timezone"`
	UserName      string    `json:"user_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	RoleName      string    `json:"role_name"`
	RoleLevel     int       `json:"role_level"`
	TenantVersion string    `json:"tenant_version"`
	UserVersion   string    `json:"user_version"`
}
