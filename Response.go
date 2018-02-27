package VkApi

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
	CallMethod *string
	CallParams *P
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

type Response struct {
	Response      *json.RawMessage `json:"response"`
	Error         ApiError         `json:"error"`
	ExecuteErrors []struct {
		Method  string `json:"method"`
		Code    int    `json:"error_code"`
		Message string `json:"error_msg"`
	} `json:"execute_errors"`
	TransportError bool
}

func (r *Response) success() bool {
	return r.Error.Code == 0
}

func (r *Response) canRetry() bool {
	switch r.Error.Code {
	case UnknownError:
		return true
	case TooManyRequests:
		return true
	case TooManyActions:
		return true
	case ServerError:
		return true
	case AdvError:
		return true
	default:
		return false
	}
}

func (r Response) Any() *AnyModel {
	m := AnyModel(*r.Response)
	return &m
}

func (r Response) Slice() ([]*AnyModel, error) {
	a := r.Any()
	return a.Slice()
}
func (r Response) FirstAny() (*AnyModel, error) {
	a := r.Any()
	s, err := a.Slice()
	if err != nil {
		return nil, err
	}
	return s[0], nil
}

func (r Response) GetString(k string) (string, error) {
	a := r.Any()
	return a.GetString(k)
}

func (r Response) GetInt(k string) (int, error) {
	a := r.Any()
	return a.GetInt(k)
}

func (r Response) GetSlice(k string) ([]*AnyModel, error) {
	a := r.Any()
	return a.GetSlice(k)
}

func (r *Response) QStringDef(k string, def string) string {
	a := r.Any()
	return a.QStringDef(k, def)
}

func (r *Response) QString(k string) *string {
	a := r.Any()
	return a.QString(k)
}

func (r *Response) QIntDef(k string, def int) int {
	a := r.Any()
	return a.QIntDef(k, def)
}

func (r *Response) QInt(k string) *int {
	a := r.Any()
	return a.QInt(k)
}

func (r Response) Unmarshal(i interface{}) error {
	return json.Unmarshal([]byte(*r.Response), i)
}

func (r Response) String() string {
	return string(*r.Response)
}
