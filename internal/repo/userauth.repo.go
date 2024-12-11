package repos

import (
	"fmt"
	"time"

	"hieupc05.github/backend-server/global"
)

type IUserAuthRepository interface {
	OTPMaker(email string, data []byte, expirationTime int64) error
}

type userAuthRepository struct{}

// AddOTP implements IUserAuthRepository
func (u *userAuthRepository) OTPMaker(email string, data []byte, expirationTime int64) error {
	key := fmt.Sprintf("usr:%s:otp", email)
	return global.Rdb.SetEx(ctx, key, data, time.Duration(expirationTime)).Err()
}

func NewUserAuthRepository() IUserAuthRepository {
	return &userAuthRepository{}
}
