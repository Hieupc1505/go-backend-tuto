package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (m *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
func (m *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	if time.Now().After(m.ExpiredAt) {
		return nil, ErrExpiredToken
	}
	return &jwt.NumericDate{Time: m.ExpiredAt}, nil
}
func (m *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: m.IssuedAt}, nil
}
func (m *Payload) GetIssuer() (string, error) {
	return "", nil
}
func (m *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: m.IssuedAt}, nil
}
func (m *Payload) GetSubject() (string, error) {
	return m.Username, nil
}
