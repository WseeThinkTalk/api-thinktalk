package errorx

type CodeError struct {
	Result ResultCode `json:"result"`
	Code   int        `json:"code"`
	Msg    string     `json:"msg"`
}

type CodeErrorResponse struct {
	Result ResultCode  `json:"result"`
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func NewCodeError(code Code) error {
	return &CodeError{Result: ResultFailure, Code: code.Value(), Msg: code.Msg()}
}

func NewDefaultError() error {
	return NewCodeError(Unknown)
}

func NewCodeMsg(code int, msg string) error {
	return &CodeError{Result: ResultFailure, Code: code, Msg: msg}
}

func NewDefaultMsg(msg string) error {
	return NewCodeMsg(int(InvalidParam), msg)
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data() *CodeErrorResponse {
	code := e.Code
	if code == OK.Value() {
		code = Unknown.Value()
	}

	return &CodeErrorResponse{
		Result: ResultFailure,
		Code:   code,
		Msg:    e.Msg,
		Data:   nil,
	}
}
