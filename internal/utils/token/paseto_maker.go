package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	util "hieupc05.github/backend-server/internal/utils"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	fmt.Println("symmetricKey: ", symmetricKey)
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly 32 characters")
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(userId int64, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userId, util.DepositorRole, duration)
	if err != nil {
		return "", payload, err
	}
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)

	if err != nil {
		return nil, ErrInvalidToken
	}
	_, err = payload.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	return payload, nil

}
