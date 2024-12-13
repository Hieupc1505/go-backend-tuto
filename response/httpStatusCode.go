package response

type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error Codes
const (
	ErrCodeSuccess = 0 // Success
	// ErrCodeParamInvalid = 20003 // Email is invalid

	// ErrInvalidToken     = 30001 // Token is invalid
	// ErrInvalidOTP       = 30002
	// ErrSendMailOtp      = 30003
	// ErrCodeUserHasExists = 50001 // User already exists
	// ErrInternalServerError = 50000
	// ErrStatusNotFound   = 50004 // Cannot find user with email
	// ErrStatusUnauthorized = 50005 // Unauthorized
	// ErrLoginFail        = 50006 // Login failed
	ErrAuthFail = 90000

	ErrSystem                       = 10001
	ErrSessionExpire                = 30001
	ErrPermissionDenied             = 40001
	ErrInvalidData                  = 50001
	ErrInvalidID                    = 50002
	ErrInvalidName                  = 50003
	ErrInvalidQuestionLevel         = 50004
	ErrInvalidQuestionAnswer        = 50005
	ErrInvalidQuestionAnswerType    = 50006
	ErrInvalidAnswerType            = 50007
	ErrInvalidQuestionType          = 50008
	ErrInvalidSubjectID             = 50009
	ErrInvalidContestTimeExam       = 50010
	ErrInvalidContestNumberQuestion = 50011
	ErrInvalidContestSubjectName    = 50012
	ErrInvalidNotEnoughAmount       = 50013
	ErrInvalidContestNotFound       = 50014
	ErrInvalidAmount                = 50015
	ErrInvalidReceiver              = 50016
	ErrInvalidTransferToOneself     = 50017
	ErrReferralNotFinished          = 50018
	ErrUserNotFound                 = 15000
	ErrDataNotFound                 = 16000
	ErrWalletNotFound               = 16001
	ErrUnauthorizedInvalidToken     = 17000
	ErrDuplicatedSelectAccountType  = 18000
	ErrInvalidContestGameID         = 19000
	ErrContestCreated               = 20001
	ErrContestRunning               = 20000
	ErrContestFinished              = 21000
	ErrContestLiveAlready           = 22000
	ErrContestLiveSubmitAlready     = 23000
	ErrProfileCheckinCompleted      = 24000
	ErrWalletBalanceNotEnough       = 25000
)

const (
	SuccessCode         = 0
	ErrSystemCode       = 1
	UnauthorizedCode    = 3
	ErrPermisionCode    = 4
	ErrInvalidCode      = 5
	ErrUserCode         = 10
	ErrDataNotFoundCode = 16
	ErrDuplicatedCode   = 18
	ErrContestStateCode = 19
	ErrWalletCode       = 20
)

// Error Messages Map
var ErrorMessages = map[int]ErrorDetail{
	ErrCodeSuccess: {SuccessCode, "Success"},
	// ErrCodeParamInvalid:          {ErrInvalidCode, "Email is invalid"},
	// ErrInvalidToken:              {ErrInvalidCode, "Token is invalid"},
	// ErrInvalidOTP:                {ErrInvalidCode, "OTP error"},
	// ErrSendMailOtp:               {ErrSystemCode, "Send mail OTP error"},
	// ErrCodeUserHasExists:         {ErrUserCode, "User already exists"},
	// ErrInternalServerError:       {ErrSystemCode, "Internal server error"},
	// ErrStatusNotFound:            {ErrUserCode, "User not found"},
	// ErrStatusUnauthorized:        {UnauthorizedCode, "Unauthorized"},
	ErrAuthFail:                     {ErrUserCode, "Auth failed"},
	ErrSystem:                       {ErrSystemCode, "system_error"},
	ErrSessionExpire:                {UnauthorizedCode, "unauthorized.session_expired"},
	ErrPermissionDenied:             {ErrPermisionCode, "permission_denied"},
	ErrInvalidData:                  {ErrInvalidCode, "invalid_data"},
	ErrInvalidID:                    {ErrInvalidCode, "invalid_data.id"},
	ErrInvalidName:                  {ErrInvalidCode, "invalid_data.name"},
	ErrInvalidQuestionLevel:         {ErrInvalidCode, "invalid_data.question.level"},
	ErrInvalidQuestionAnswer:        {ErrInvalidCode, "invalid_data.question.answer"},
	ErrInvalidQuestionAnswerType:    {ErrInvalidCode, "invalid_data.question.answer_type"},
	ErrInvalidAnswerType:            {ErrInvalidCode, "Invalid answer_type's question"},
	ErrInvalidQuestionType:          {ErrInvalidCode, "invalid_data.question.type"},
	ErrInvalidSubjectID:             {ErrInvalidCode, "invalid_data.subject_id"},
	ErrInvalidContestTimeExam:       {ErrInvalidCode, "invalid_data.contest.time_exam"},
	ErrInvalidContestNumberQuestion: {ErrInvalidCode, "invalid_data.contest.number_question"},
	ErrInvalidContestSubjectName:    {ErrInvalidCode, "invalid_data.contest.subject_name"},
	ErrInvalidNotEnoughAmount:       {ErrInvalidCode, "invalid_data.not_enough_amount"},
	ErrInvalidContestNotFound:       {ErrInvalidCode, "invalid_data.contest.notfound"},
	ErrInvalidAmount:                {ErrInvalidCode, "invalid_data.amount"},
	ErrInvalidReceiver:              {ErrInvalidCode, "invalid_data.receiver"},
	ErrInvalidTransferToOneself:     {ErrInvalidCode, "invalid_data.transfer_to_oneself"},
	ErrReferralNotFinished:          {ErrUserCode, "referral.not_finished"},
	ErrUserNotFound:                 {ErrUserCode, "user.not_found"},
	ErrDataNotFound:                 {ErrDataNotFoundCode, "data_not_found"},
	ErrWalletNotFound:               {ErrWalletCode, "wallet.not_found"},
	ErrUnauthorizedInvalidToken:     {UnauthorizedCode, "unauthorized.missing_or_invalid_token"},
	ErrDuplicatedSelectAccountType:  {ErrDuplicatedCode, "duplicated.only_select_account_type_once"},
	ErrInvalidContestGameID:         {ErrInvalidCode, "invalid_data.contest.game_id"},
	ErrContestLiveSubmitAlready:     {ErrContestStateCode, "contest.live.submit_already"},
	ErrWalletBalanceNotEnough:       {ErrWalletCode, "wallet.balance.not_enough"},
	ErrContestCreated:               {0, ""},
}
