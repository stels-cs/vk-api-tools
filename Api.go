package Vk

import (
	"encoding/json"
	"time"
	"log"
)

type Method string
type Params map[string]string

func (p *Params) toString() string {
	str := ""
	for k, v := range *p {
		if str != "" {
			str += " "
		}
		str += k + "=" + v
	}
	return str
}

func (p *Params) toMap() map[string][]string {
	out := map[string][]string{}
	for k, v := range *p {
		out[k] = []string{v}
	}
	return out
}


func (p *Params) toJson() ([]byte, error) {
	return json.Marshal(p)
}

type CaptchaListener interface {
	OnCaptcha(image string)
}

type Api struct {
	maxRetryCount int
	token         AccessToken
	transport     Transport
	version       string
	Message       *Message
	Users         *Users
	Groups         *Groups
	logger        *log.Logger
	Captcha       chan string
	captchaSid string
	captchaKey string
	captchaListener CaptchaListener
}

func GetApi(token AccessToken, transport Transport, logger *log.Logger) *Api {
	m := Message{}
	u := Users{}
	g := Groups{}
	api := &Api{
		30,
		token,
		transport,
		"5.69",
		&m,
		&u,
		&g,
		logger,
		make(chan string, 10),
		"",
		"",
		nil,
	}
	m.api = api
	u.api = api
	g.api = api
	return api
}

func (api *Api) run(method Method, params Params, retryCount int) (ApiResponse, error) {
	response := ApiResponse{}

	params["access_token"] = api.token.Token
	params["v"] = api.version
	if api.captchaSid != "" && api.captchaKey != "" {
		params["captcha_sid"] = api.captchaSid
		params["captcha_key"] = api.captchaKey
	}
	data, err, external := api.transport.call(method, params)
	if err != nil {
		if retryCount < api.maxRetryCount {
			return api.retryCall(method, params, retryCount, err)
		} else {
			return response, err
		}
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		if retryCount < api.maxRetryCount {
			return api.retryCall(method, params, retryCount, err)
		} else {
			return response, &TransportBadResponse{
				method,
				params,
				data,
				external,
			}
		}
	}

	api.captchaSid = ""
	api.captchaKey = ""

	if response.Error.Code == ApiErrorCaptcha {
		api.onCaptchaFired(response.Error.CaptchaSid, response.Error.CaptchaImg)
		response.Error.CallMethod = &method
		response.Error.CallParams = &params
		return response, &response.Error
	}

	if response.canRetry() && retryCount < api.maxRetryCount {
		response.Error.CallMethod = &method
		response.Error.CallParams = &params
		return api.retryCall(method, params, retryCount, &response.Error)
	}

	if response.success() {
		return response, nil
	} else {
		return response, &response.Error
	}
}

func (api *Api) onCaptchaFired(captchaSid string, captchaImg string) {
	if api.captchaSid == "" {
		api.captchaSid = captchaSid
		api.captchaKey = ""
		if api.captchaListener != nil {
			api.captchaListener.OnCaptcha(captchaImg)
		}
	}
}

func (api *Api) SetCaptchaListener(listener CaptchaListener) {
	api.captchaListener = listener
}

func (api *Api) SetCaptchaKey(key string)  {
	api.captchaKey = key
}

func (api *Api) retryCall(method Method, params Params, retryCount int, err error) (ApiResponse, error) {
	time.Sleep(time.Second)
	if api.logger != nil {
		if err != nil  {
			api.logger.Println(PrintError(err))
		}
	}
	res, err := api.run(method, params, retryCount+1)
	return res, err
}

func (api *Api) BlindExecute(code string) error {
	_, err := api.Execute(code)
	return err
}

func (api *Api) Execute(code string) (ApiResponse, error) {
	return api.run("execute", map[string]string{"code": code}, 0)
}

func (api *Api) Call( method ApiMethod ) (ApiResponse, error) {
	return api.run(method.name,  method.params, 0)
}