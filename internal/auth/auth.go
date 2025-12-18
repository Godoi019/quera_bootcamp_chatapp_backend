package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	symmetricKey    paseto.V4SymmetricKey
	tokenExpiration time.Duration
}

type TokenPayload struct {
	UserID   int       `json:"user_id"`
	Username string    `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpireAt time.Time `json:"expire_at"`
}

func NewAuthService(pasetoKey string, expirationHours int) *Service {
	// Create a V4 symmetric key from the provided key string
	// For production, use paseto.NewV4SymmetricKey() to generate a secure key
	key, err := paseto.V4SymmetricKeyFromBytes([]byte(pasetoKey))
	if err != nil {
		// If the key is invalid, generate a new one (should log warning in production)
		key = paseto.NewV4SymmetricKey()
	}

	return &Service{
		symmetricKey:    key,
		tokenExpiration: time.Duration(expirationHours) * time.Hour,
	}
}

func (s *Service) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (s *Service) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) CreateToken(userID int, username string) (string, error) {
	now := time.Now()
	payload := TokenPayload{
		UserID:   userID,
		Username: username,
		IssuedAt: now,
		ExpireAt: now.Add(s.tokenExpiration),
	}

	// Marshal payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create PASETO v4 token
	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(payload.ExpireAt)
	token.SetString("data", string(payloadJSON))

	// Encrypt the token
	encrypted := token.V4Encrypt(s.symmetricKey, nil)

	return encrypted, nil
}

func (s *Service) VerifyToken(tokenString string) (*TokenPayload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	// Parse and verify the token
	token, err := parser.ParseV4Local(s.symmetricKey, tokenString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	// Extract the payload
	dataStr, err := token.GetString("data")
	if err != nil {
		return nil, fmt.Errorf("failed to get token data: %w", err)
	}

	// Unmarshal the payload
	var payload TokenPayload
	if err := json.Unmarshal([]byte(dataStr), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &payload, nil
}
