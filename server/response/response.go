package response

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "成功",
		Data:    data,
	}
}

func Error(msg string) Response {
	return Response{
		Code:    500,
		Message: msg,
	}
}
