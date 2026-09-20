package errorx

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

type ResponseSuccessBean struct {
	Result ResultCode  `json:"result"`
	Code   Code        `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func Success(data interface{}) *ResponseSuccessBean {
	return &ResponseSuccessBean{
		Result: ResultSuccess,
		Code:   OK,
		Msg:    FindMsg(OK),
		Data:   data,
	}
}

type ResponseErrorBean struct {
	Result ResultCode  `json:"result"`
	Code   Code        `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func DefinedResponse(result ResultCode, errCode Code, errMsg string, data interface{}) *ResponseErrorBean {
	if FindMsg(errCode) != "" && errMsg == "" {
		errMsg = FindMsg(errCode)
	}

	return &ResponseErrorBean{
		Result: result,
		Code:   errCode,
		Msg:    errMsg,
		Data:   data,
	}
}

func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if resp == nil {
		resp = make(map[string]interface{})
	}

	// success
	if err == nil {
		httpx.WriteJson(w, http.StatusOK, Success(resp))
		return
	}

	// error
	var (
		errCode = Unknown
		errMsg  string
		result  = ResultFailure
	)

	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*CodeError); ok {
		errCode = Code(e.Code)
		errMsg = e.Msg
		result = e.Result
	} else if st, ok := status.FromError(err); ok {
		errCode = Code(st.Code())
		errMsg = st.Message()
	} else {
		errMsg = err.Error()
	}

	if errMsg == "" {
		errMsg = FindMsg(Unknown)
	}

	httpx.WriteJson(w, http.StatusOK, DefinedResponse(result, errCode, errMsg, resp))
}
