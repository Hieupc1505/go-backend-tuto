package api

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "hieupc05.github/backend-server/db/mock"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/services"
	ucrypto "hieupc05.github/backend-server/internal/utils/crypto"
	passExt "hieupc05.github/backend-server/internal/utils/password"
	"hieupc05.github/backend-server/internal/utils/random"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = random.RandomString(6)
	hashedPassword, err := passExt.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		HashedPassword: hashedPassword,
		Email:          random.RandomEmail(),
	}
	return

}

func TestRegisterUser(t *testing.T) {
	// user := randomAccount()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    "hoanghieuss3x@gmail.com",
				"password": random.RandomString(8),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("hoanghieuss3x@gmail.com")).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "EmailExists",
			body: gin.H{
				"email":    "sKf3o@example.com",
				"password": random.RandomString(8),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			router := newTestServer(store)
			recorder := httptest.NewRecorder()

			url := "/v1/2024/user/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder)
		})
	}

}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := passExt.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func eqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func SetupCreateUser(ctx context.Context, email string, password string, defaultOtp int) {
	hashMail := ucrypto.GetHash(email)
	hashPass, _ := passExt.HashPassword(password)
	data, _ := json.Marshal(services.RegisterDataHash{
		OTP:      strconv.Itoa(int(defaultOtp)),
		Password: hashPass,
	})
	global.Rdb.SetEx(ctx, fmt.Sprintf("usr:%s:otp", hashMail), []byte(data), time.Minute)
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)
	defaultOtp := 123456
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setup         func(c context.Context)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email": user.Email,
				"otp":   defaultOtp,
			},
			setup: func(ctx context.Context) {
				SetupCreateUser(ctx, user.Email, password, defaultOtp)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email: user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(arg, password)).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidEmailAndMailNotExist",
			body: gin.H{
				"email": user.Email,
				"otp":   123456,
			},
			setup: func(ctx context.Context) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "IncorrectOtp",
			body: gin.H{
				"email": user.Email,
				"otp":   123457,
			},
			setup: func(ctx context.Context) {
				SetupCreateUser(ctx, user.Email, password, defaultOtp)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			router := newTestServer(store)
			recorder := httptest.NewRecorder()

			url := "/v1/2024/user/verify_otp"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			// Create a context and call the setup function
			ctx := context.Background()
			tc.setup(ctx)

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: gin.H{
				"email":    "notfound@gmail.com",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"email":    user.Email,
				"password": "IncorrectPassword",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"email":    "invalidemail",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			router := newTestServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/v1/2024/user/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}
