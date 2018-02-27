# ApiError

```go
type ApiError struct {
	Code       int               `json:"error_code"`
	Message    string            `json:"error_msg"`
	Params     []RequestedParams `json:"request_params"`
	CaptchaSid string            `json:"captcha_sid"`
	CaptchaImg string            `json:"captcha_img"`
	CallMethod *string // Ссылка на метод который был вызван
	CallParams *VkApi.P // Ссылка на парамтеры которые были переданы
}
```

[Назад](../README.md)