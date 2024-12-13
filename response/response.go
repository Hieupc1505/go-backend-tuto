package response

type Response struct {
	Errors interface{} `json:"e"`
	Data   interface{} `json:"d,omitempty"`
}

type MsgResponse struct {
	Msg string `json:"msg"`
}

type EmptyObj struct{}

func SuccessResponse(code int, data interface{}) Response {
	res := Response{
		Data:   data,
		Errors: 0,
	}
	return res
}

func ErrorResponse(code int) Response {
	msg := ErrorMessages[code]
	res := Response{
		Errors: msg.Code,
		Data: MsgResponse{
			Msg: msg.Message,
		},
	}
	return res
}
