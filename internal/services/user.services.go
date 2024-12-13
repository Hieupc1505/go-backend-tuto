package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
	repos "hieupc05.github/backend-server/internal/repo"
	ucrypto "hieupc05.github/backend-server/internal/utils/crypto"
	"hieupc05.github/backend-server/internal/utils/password"
	passExt "hieupc05.github/backend-server/internal/utils/password"
	"hieupc05.github/backend-server/internal/utils/random"
	"hieupc05.github/backend-server/internal/utils/sendto"
	"hieupc05.github/backend-server/internal/utils/token"
	"hieupc05.github/backend-server/response"
)

// type Services struct {
// 	info *repos.Repo
// }

// func NewService() *Services {

// 	return &Services{
// 		info: repos.UserRepo(),
// 	}
// }

// func (c *Services) GetUserInfoService() string {
// 	return c.info.GetUserInfo()
// }

type IUserServices interface {
	Register(email string, password string) (r response.Response, httpStatus int)
	CreateUser(ctx *gin.Context, otp int, email string) (r response.Response, httpStatus int)
	Login(ctx *gin.Context, email string, password string) (r response.Response, httpStatus int)
}

type RegisterDataHash struct {
	OTP      string `json:"otp"`
	Password string `json:"password"`
}

type userServices struct {
	userRepo     repos.IUserRepository
	userAuthRepo repos.IUserAuthRepository
	tokenMaker   token.Maker
}

type RegisterResult struct {
	ErrorCode  int         // Mã lỗi
	HttpStatus int         // Mã HTTP
	Data       interface{} // Dữ liệu trả về (nếu có)
}

type loginUserResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  db.User   `json:"user"`
}

func (us *userServices) Register(email string, password string) (_ response.Response, httpStatus int) {
	// 1. Check email exist
	if _, err := us.userRepo.GetUserByEmail(email); err == nil {
		return response.ErrorResponse(response.ErrAuthFail), http.StatusConflict
	}
	fmt.Println("get user by email + 1")
	// 2. new OTP
	otp := random.RandomInt(100000, 999999)
	if global.Config.Server.Mode == "dev" {
		otp = 123456
	}

	hasPass, err := passExt.HashPassword(password)
	if err != nil {
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	data, err := json.Marshal(RegisterDataHash{
		OTP:      strconv.Itoa(int(otp)),
		Password: hasPass,
	})
	if err != nil {
		return response.ErrorResponse(response.ErrInvalidData), http.StatusInternalServerError
	}

	hashEmail := ucrypto.GetHash(email)

	// 3. Save OTP in Redis with expiration time
	err = us.userAuthRepo.OTPMaker(hashEmail, data, int64(10*time.Minute))
	if err != nil {
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	// 4. Send Email OTP to user
	err = sendto.SendTemplateEmailOtp([]string{email}, "hoanghieuss3344@gmail.com", "otp-auth.html", map[string]interface{}{
		"otp": strconv.Itoa(int(otp)),
	})

	if err != nil {
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	// Success
	return response.SuccessResponse(response.ErrCodeSuccess, nil), http.StatusOK
}

func (us *userServices) CreateUser(ctx *gin.Context, otp int, email string) (_ response.Response, httpStatus int) {
	// Hash email to use as Redis key
	hashMail := ucrypto.GetHash(email)

	// 1. Get and parse OTP data from Redis
	var otpData RegisterDataHash
	otpRedis, err := global.Rdb.Get(ctx, fmt.Sprintf("usr:%s:otp", hashMail)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return response.ErrorResponse(response.ErrSystem), http.StatusBadRequest
		}
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	// Parse the OTP data
	if err := json.Unmarshal([]byte(otpRedis), &otpData); err != nil {
		return response.ErrorResponse(response.ErrInvalidData), http.StatusBadRequest
	}

	// Convert OTP from Redis and validate
	otpCnv, err := strconv.Atoi(otpData.OTP)
	if err != nil || otp != otpCnv {
		return response.ErrorResponse(response.ErrInvalidData), http.StatusBadRequest
	}

	// 2. Save user to the database
	if _, err := global.PgDb.CreateUser(ctx, db.CreateUserParams{
		Email:          email,
		HashedPassword: otpData.Password,
	}); err != nil {
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	// 3. Delete OTP from Redis
	if err := global.Rdb.Del(ctx, fmt.Sprintf("usr:%s:otp", hashMail)).Err(); err != nil {
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	// 4. Return success response
	return response.SuccessResponse(response.ErrCodeSuccess, nil), http.StatusOK
}

func (us *userServices) Login(ctx *gin.Context, email string, purpose string) (_ response.Response, httpStatus int) {

	user, err := us.userRepo.GetUserByEmail(email)
	if err != nil {
		if err == db.ErrRecordNotFound {
			return response.ErrorResponse(response.ErrUserNotFound), http.StatusBadRequest
		}
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	err = password.CheckPassword(purpose, user.HashedPassword)
	if err != nil {
		return response.ErrorResponse(response.ErrAuthFail), http.StatusUnauthorized
	}

	accessToken, accessPayload, err := us.tokenMaker.CreateToken(user.ID, user.Role, global.Config.Token.AccessTokenDuration)
	if err != nil {
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	refreshToken, refreshPayload, err := us.tokenMaker.CreateToken(user.ID, user.Role, global.Config.Token.RefreshTokenDuration)
	if err != nil {
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	session, err := global.PgDb.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return response.ErrorResponse(response.ErrSystem), http.StatusInternalServerError
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  user,
	}

	return response.SuccessResponse(response.ErrCodeSuccess, rsp), http.StatusOK
}

func NewUserServices(userRepo repos.IUserRepository, userAuthRepo repos.IUserAuthRepository, tokenMaker token.Maker) IUserServices {
	return &userServices{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
		tokenMaker:   tokenMaker,
	}
}
