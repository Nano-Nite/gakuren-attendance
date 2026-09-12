package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"

	"strings"

	"gakuren-system.com/pkg/db"
	"gakuren-system.com/pkg/helper"
	"gakuren-system.com/pkg/model"
)

var JWTService *model.JWTService

func NewJWTService() error {
	privateKey, err := helper.ParsePrivateKey(os.Getenv("RSA_PRIVATE_KEY"))
	if err != nil {
		return fmt.Errorf("parse JWT private key: %w", err)
	}

	service := new(model.JWTService)
	service.PrivateKey = privateKey

	JWTService = service

	log.Println("JWT Service Started")
	return nil
}

func AuthenticateBearerToken(authHeader string) (string, error) {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("authorization header must be in format: Bearer <token>")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("bearer token is required")
	}

	expectedToken := strings.TrimSpace(os.Getenv("BEARER_TOKEN"))
	if expectedToken == "" {
		expectedToken = "dev-token"
	}

	if token != expectedToken {
		return "", errors.New("invalid bearer token")
	}

	return token, nil
}

func GetPublicKey() (string, error) {
	publicKeyEnv := strings.TrimSpace(os.Getenv("RSA_PUBLIC_KEY"))
	if publicKeyEnv == "" {
		return "", errors.New("public key is not set in environment variables")
	}

	publicKeyPem, err := decodePEM(publicKeyEnv)
	if err != nil {
		publicKeyPem = publicKeyEnv
	}

	publicKeyPem = strings.TrimSpace(publicKeyPem)
	publicKeyPem = strings.ReplaceAll(publicKeyPem, `\n`, "\n")

	if !strings.Contains(publicKeyPem, "-----BEGIN PUBLIC KEY-----") && !strings.Contains(publicKeyPem, "-----BEGIN RSA PUBLIC KEY-----") {
		return "", errors.New("RSA_PUBLIC_KEY does not contain a valid PEM public key")
	}

	return publicKeyPem, nil
}

func decodePEM(src string) (string, error) {
	tmp := strings.TrimSpace(src)
	if tmp == "" {
		return "", errors.New("empty base64 input")
	}

	decoded, err := helper.DecodeB64Bytes(tmp)
	if err == nil {
		return string(decoded), nil
	}

	decoded, err = base64.URLEncoding.DecodeString(tmp)
	if err == nil {
		return string(decoded), nil
	}

	decoded, err = base64.RawStdEncoding.DecodeString(tmp)
	if err == nil {
		return string(decoded), nil
	}

	decoded, err = base64.RawURLEncoding.DecodeString(tmp)
	if err == nil {
		return string(decoded), nil
	}

	return "", err
}

func GetLoginDataByEmail(email string) (*model.LoginData, error) {
	loginData, err := db.GetSingleDataByQuery[model.LoginData](`
	select
		u.uuid, t."uuid" as tenant_uuid,
		t.name tenant_name,
		u.school_uuid as school_uuid,
		sc.code,
		t.timezone,
		t.version tenant_version,
		u.name user_name,
		u.email,
		u.phone,
		u.address,
		u.version user_version,
		r.name role_name,
		r.level role_level
	from user_sch.tenant t
	join user_sch.user u on t.uuid = u.tenant_uuid
	join user_sch.role r on u.role_uuid = r.uuid
	join school_sch.school sc on u.school_uuid = sc.uuid
	where u.email = $1
	`, email)
	if err != nil {
		return nil, err
	}
	return loginData, nil
}
