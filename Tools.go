package VkApi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
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

// Так можно проверить что err который вернулся при запросе к API – это ошибка типа TransportError
// Ошибки такого типа возникают из-за проблем с сетью или когда ВК недоступен
func IsTransportError(err interface{}) bool {
	if _, ok := err.(TransportError); ok {
		return true
	}
	if _, ok := err.(*TransportError); ok {
		return true
	}
	return false
}

// Так можно проверить что err который вернулся при запросе к API – это ошибка типа ApiError
// Ошибки такого типа возникают из-за проблем с запросом, например нет прав или слишком много запросов в секунду
func IsApiError(err interface{}) bool {
	if _, ok := err.(ApiError); ok {
		return true
	}
	if _, ok := err.(*ApiError); ok {
		return true
	}
	return false
}

// Так можно проверить что err который вернулся при запросе к API – это ошибка типа ApiError и код ошибки - 14
// ВК просит нас ввести капчу
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

func CalculateSignature(query map[string][]string, pattern string, secret string) string {
	if secret == "" {
		return "EMPTY_SECRET_" + strconv.Itoa(rand.Int())
	}

	buff := ""

	index := make([]int, 0, len(query))
	keys := make([]string, 0, len(query))
	keyIndex := make(map[int]string, 0)

	for key := range query {
		i := strings.Index(pattern, key+"=")
		keys = append(keys, key)
		keyIndex[i] = key
		index = append(index, i)
	}
	sort.Ints(index)

	for _, i := range index {
		key := keyIndex[i]
		if key == "hash" || key == "sign" || key == "api_result" {
			continue
		}
		payload := ""
		if len(query[key]) != 0 {
			payload = query[key][0]
		}
		if key == "ad_info" {
			payload = strings.Replace(payload, " ", "+", -1)
		}
		buff += payload
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(buff))
	sign := mac.Sum(nil)
	return hex.EncodeToString(sign)
}

// Проверка что подпись запроса верна https://vk.com/dev/community_apps_docs
func IsCorrectRequest(query string, secret string) bool {

	v, err := url.ParseQuery(query)
	if err != nil {
		return false
	}

	s := v["sign"]

	if len(s) != 1 {
		return false
	}

	_s := CalculateSignature(v, query, secret)

	return _s == s[0]
}
