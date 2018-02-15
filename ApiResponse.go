package Vk

import (
	"encoding/json"
	"strconv"
)

type RequestedParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ApiError struct {
	Code       int               `json:"error_code"`
	Message    string            `json:"error_msg"`
	Params     []RequestedParams `json:"request_params"`
	CaptchaSid string            `json:"captcha_sid"`
	CaptchaImg string            `json:"captcha_img"`
	CallMethod *Method
	CallParams *Params
}

func (e *ApiError) Error() string {
	str := strconv.Itoa(e.Code) + " " + e.Message
	if len(e.Params) > 0 {
		str += " " + requestParamsToString(e.Params)
	}
	if e.CallMethod != nil {
		str += " run:" + string(*e.CallMethod)
	}
	if e.CallParams != nil {
		str += " " + e.CallParams.toString()
	}

	return str
}

type ApiResponse struct {
	Response *json.RawMessage `json:"response"`
	Error    ApiError         `json:"error"`
	ExecuteErrors []struct {
		Method  string `json:"method"`
		Code    int    `json:"error_code"`
		Message string `json:"error_msg"`
	} `json:"execute_errors"`
}

func (r *ApiResponse) success() bool {
	return r.Error.Code == 0
}

func (r *ApiResponse) canRetry() bool {
	switch r.Error.Code {
	case ApiUnknowmError:
		return true
	case ApiTooManyRequests:
		return true
	case ApiTooManyActions:
		return true
	case ApiSearverError:
		return true
	case ApiError603:
		return true
	default:
		return false
	}
}
