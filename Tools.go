package Vk

import (
	"fmt"
	"strconv"
	"encoding/json"
)

func iif(cond bool, ok string, bad string) string  {
	if cond {
		return ok
	} else {
		return bad
	}
}

func requestParamsToString(params []RequestedParams) string {
	str := ""
	for _, v := range params {
		if str != "" {str += " "}
		str += v.Key + "=" + v.Value
	}
	return str
}


func callToString(method Method, params Params) string  {
	return string(method) + " " + params.toString()
}

func PrintError(err error) string  {
	if e, ok := err.(*AuthError); ok {
		return fmt.Sprintf("[AuthError] %s", e.Error())
	} else if e, ok := err.(*ApiError); ok {
		if e.CaptchaSid != "" {
			return fmt.Sprintf("[CaptchaError] key %s image %s %s", e.CaptchaSid, e.CaptchaImg, e.Error())
		} else {
			return fmt.Sprintf("[ApiError] %s", e.Error())
		}
	} else if e, ok := err.(*TransportBadResponse); ok {
		return fmt.Sprintf("[TransportBadResponse] %s", e.Error())
	} else {
		return fmt.Sprintf("[Error] %s", err.Error())
	}
}

func intToString( items []int ) string {
	str := ""
	for _,v:=range items {
		if str != "" {
			str += ","
		}
		str += strconv.Itoa(v)
	}
	return str
}

func isBoolAndFalse( raw *json.RawMessage ) bool {
	var d interface{}
	err := json.Unmarshal( *raw, &d )
	if err != nil {
		return false
	}
	if b, ok := d.(bool); ok && b == false {
		return true
	}
	return false
}