package response

type Response struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type EmptyObj struct{}

func SuccessResponse(code int, message string, data interface{}) Response {
	res := Response{
		Message: message,
		Data:    data,
		Code:    0,
		Errors:  nil,
	}
	return res
}

func ErrorResponse(code int, message string, err string) Response {
	res := Response{
		Code:    code,
		Message: msg[code],
		Errors:  err,
		Data:    nil,
	}
	return res
}
