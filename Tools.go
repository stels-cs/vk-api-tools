package VkApi

import (
	"encoding/json"
	"fmt"
)

var defaultApi *Api

func getDefaultApi() *Api {
	if defaultApi == nil {
		defaultApi = CreateApi("", "5.71", GetHttpTransport(), 30)
	}
	return defaultApi
}

func requestParamsToString(params []RequestedParams) string {
	str := ""
	for _, v := range params {
		if str != "" {
			str += " "
		}
		str += v.Key + "=" + v.Value
	}
	return str
}

func callToString(method string, params P) string {
	return string(method) + " " + params.toString()
}

func printError(err error) string {
	if e, ok := err.(*ApiError); ok {
		if e.CaptchaSid != "" {
			return fmt.Sprintf("[CaptchaError] key %s image %s %s", e.CaptchaSid, e.CaptchaImg, e.Error())
		} else {
			return fmt.Sprintf("[ApiError] %s", e.Error())
		}
	} else if e, ok := err.(*TransportError); ok {
		return fmt.Sprintf("[TransportError] %s", e.Error())
	} else {
		return fmt.Sprintf("[Error] %s", err.Error())
	}
}

func isBoolAndFalse(raw *json.RawMessage) bool {
	var d interface{}
	err := json.Unmarshal(*raw, &d)
	if err != nil {
		return false
	}
	if b, ok := d.(bool); ok && b == false {
		return true
	}
	return false
}

func Call(method string, params P) (Response, error) {
	api := getDefaultApi()
	return api.Call(string(method), params)
}

// Выполняет запрос к API ВКонтакте, в случае любых серевых ошибок или кодов ошибок (1, 9, 6, 9, 10, 603) повторяет запрос до 30 раз.
// 	users := make([]struct{
// 		Id        int    `json:"id"`
// 		FirstName string `json:"first_name"`
// 		LastName  string `json:"last_name"`
// 	}, 0)
// 	err := VkApi.Exec("users.get", VkApi.P{"user_ids": "2050,andrew"}, &users)
func Exec(method string, params P, s interface{}) error {
	api := getDefaultApi()
	return api.Exec(string(method), params, s)
}

func IsTransportError(err interface{}) bool {
	if _, ok := err.(TransportError); ok {
		return true
	}
	if _, ok := err.(*TransportError); ok {
		return true
	}
	return false
}

func IsApiError(err interface{}) bool {
	if _, ok := err.(ApiError); ok {
		return true
	}
	if _, ok := err.(*ApiError); ok {
		return true
	}
	return false
}

func IsCaptchaError(err interface{}) bool {
	if IsApiError(err) {
		e := CastToApiError(err)
		return e.Code == CaptchaError
	}
	return false
}

func CastToApiError(err interface{}) *ApiError {
	if e, ok := err.(ApiError); ok {
		return &e
	}
	if e, ok := err.(*ApiError); ok {
		return e
	}
	panic("Cant cast error to *ApiError")
}

func CastToTransportError(err interface{}) *TransportError {
	if e, ok := err.(TransportError); ok {
		return &e
	}
	if e, ok := err.(*TransportError); ok {
		return e
	}
	panic("Cant cast error to *ApiError")
}
