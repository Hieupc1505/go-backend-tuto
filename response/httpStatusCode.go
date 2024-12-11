package response

const (
	ErrCodeSuccess      = 20001 //Success
	ErrCodeParamInvalid = 20003 //Email is invalid

	ErrInvalidToke = 30001 //token is invalid
	ErrInvalidOTP  = 30002
	ErrSendMailOtp = 30003
	//REgister Code
	ErrCodeUserHasExists    = 50001 //user has exists
	ErrInternalServerlError = 50000
	ErrStatusNotFound       = 50004 //Can not find user with email
	ErrStatusUnauthorized   = 50005 //Unauthorized
	ErrLoginFail            = 50006 //fail login
)

// message
var msg = map[int]string{
	ErrCodeSuccess:      "Success",
	ErrCodeParamInvalid: "Email is invalid",
	ErrInvalidToke:      "token is invalid",
	ErrInvalidOTP:       "OTP error",
	ErrSendMailOtp:      "send mail error",
	ErrLoginFail:        "Login Fail",

	ErrCodeUserHasExists:    "user has exists",
	ErrInternalServerlError: "Internal server error",
	ErrStatusNotFound:       "User not found",
	ErrStatusUnauthorized:   "Unauthorized",
}
