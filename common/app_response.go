package common

const (
	CodeOk               = 0
	CodeFail             = 1
	CodeUserExisted      = 10
	CodeMethodNotSupport = 11
)

var message = map[int]string{
	CodeOk:               "success",
	CodeFail:             "fail",
	CodeUserExisted:      "username đã tồn tại",
	CodeMethodNotSupport: "phương thức không được hỗ trợ",
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(code int, message string, data interface{}) *response {
	return &response{Code: code, Message: message, Data: data}
}

func NewSuccessResponse(data interface{}) *response {
	return NewResponse(CodeOk, GetMessageFromCode(CodeOk), data)
}

func GetMessageFromCode(code int) string {
	return message[code]
}
