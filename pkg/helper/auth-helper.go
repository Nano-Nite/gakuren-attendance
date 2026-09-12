package helper

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"gakuren-system.com/pkg/db"
	"gakuren-system.com/pkg/model"
	"github.com/google/uuid"
	"github.com/youmark/pkcs8"
)

func DecodeRSA(data string) ([]byte, error) {
	// Decode Base64 ciphertext
	ciphertext, err := DecodeB64Bytes(data)
	if err != nil {
		return nil, fmt.Errorf("invalid RSA ciphertext Base64: %w", err)
	}

	// Load private key
	rsaPrivateKey, err := ParsePrivateKey(os.Getenv("RSA_PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}

	// RSA-2048 / OAEP / SHA-256
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("RSA OAEP SHA-256 decrypt failed: %w", err)
	}

	return plaintext, nil
}

func EncodeRSA(data string) (string, error) {
	// Load RSA public key
	publicKey, err := ParsePublicKey(os.Getenv("RSA_PUBLIC_KEY"))
	if err != nil {
		return "", fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	// UTF-8 bytes
	plaintext := []byte(data)

	// RSA-OAEP SHA-256
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, nil)
	if err != nil {
		return "", fmt.Errorf("RSA encryption failed: %w", err)
	}

	// Ciphertext → Base64
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

func ParsePublicKey(base64Key string) (*rsa.PublicKey, error) {
	key := strings.TrimSpace(base64Key)
	if key == "" {
		return nil, fmt.Errorf("public key is empty")
	}

	pemBytes := []byte(key)
	if !strings.Contains(key, "BEGIN") {
		decoded, err := DecodeB64Bytes(key)
		if err == nil {
			pemBytes = decoded
		}
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not RSA public key")
	}

	return rsaPublicKey, nil
}

func ParsePrivateKey(base64Key string) (*rsa.PrivateKey, error) {
	key := strings.TrimSpace(base64Key)
	if key == "" {
		return nil, fmt.Errorf("private key is empty")
	}

	pemBytes := []byte(key)
	if !strings.Contains(key, "BEGIN") {
		decoded, err := base64.StdEncoding.DecodeString(key)
		if err == nil {
			pemBytes = decoded
		}
	}

	block, rest := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM")
	}

	if len(rest) > 0 {
		fmt.Println("warning: extra data after PEM")
	}

	if block.Type != "ENCRYPTED PRIVATE KEY" && block.Type != "PRIVATE KEY" && block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM type: %s", block.Type)
	}

	var privateKey interface{}
	switch block.Type {
	case "ENCRYPTED PRIVATE KEY":
		parsedKey, err := pkcs8.ParsePKCS8PrivateKey(block.Bytes, []byte(os.Getenv("RSA_PASSPHARSE")))
		if err != nil {
			return nil, fmt.Errorf("PKCS#8 decrypt failed: %w", err)
		}
		privateKey = parsedKey
	case "PRIVATE KEY", "RSA PRIVATE KEY":
		parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err == nil {
			privateKey = parsedKey
		} else {
			parsedKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("parse private key: %w", err)
			}
			privateKey = parsedKeyInterface
		}
	}

	if privateKey == nil {
		return nil, fmt.Errorf("unsupported private key format: %s", block.Type)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		if parsed, ok := privateKey.(*rsa.PrivateKey); ok {
			rsaPrivateKey = parsed
		} else {
			return nil, fmt.Errorf("expected RSA private key, got %T", privateKey)
		}
	}

	if err := rsaPrivateKey.Validate(); err != nil {
		return nil, fmt.Errorf("RSA validation failed: %w", err)
	}

	publicRSA, err := ParsePublicKey(os.Getenv("RSA_PUBLIC_KEY"))
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	if rsaPrivateKey.N.Cmp(publicRSA.N) != 0 {
		return nil, fmt.Errorf("public/private key pair does not match")
	}

	// log.Println("Public/Private key pair: MATCH")

	return rsaPrivateKey, nil
}

func GenerateSecureKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	os.Setenv("JWT_REFRESH_SECURE_KEY", hex.EncodeToString(bytes))
	return hex.EncodeToString(bytes), nil
}

func GetUserPermission(userUUID string, permission string) (bool, error) {
	query := `
	select distinct p.code from user_sch."user" u 
	join user_sch.role r on u.role_uuid = r.uuid
	join user_sch.role_permission rp on r.uuid = rp.role_uuid
	join user_sch.permission p on rp.permission_uuid = p.uuid
	where u.uuid = $1 and p.code = $2
	limit 1`

	selectedPermission, err := db.GetSingleDataByQuery[model.GetPermission](query, userUUID, permission)
	if err != nil {
		return false, err
	}

	if selectedPermission == nil {
		return false, nil
	}

	return true, nil
}
func GetUserStatus(statusCode string) (*uuid.UUID, error) {
	var statusUUID uuid.UUID
	err := db.Conn.QueryRow(context.Background(), `select uuid from public.status where lower(category) = 'student_enrollment' and lower(code) = lower($1)`, statusCode).Scan(&statusUUID)
	if err != nil {
		return nil, err
	}
	return &statusUUID, nil
}

func ApprovalBypass(userUUID string) (bool, error) {
	query := `
	select distinct p.code from user_sch."user" u 
	join user_sch.role r on u.role_uuid = r.uuid
	join user_sch.role_permission rp on r.uuid = rp.role_uuid
	join user_sch.permission p on rp.permission_uuid = p.uuid
	where u.uuid = $1 and p.code = $2
	limit 1`

	selectedPermission, err := db.GetSingleDataByQuery[model.GetPermission](query, userUUID, APPROVAL_BYPASS)
	if err != nil {
		return false, err
	}

	if selectedPermission == nil {
		return false, nil
	}

	return true, nil
}
