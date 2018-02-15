package Vk

import (
	"net/http"
	"time"
	"io/ioutil"
)


type HttpTransport struct {
	client *http.Client
	endpoint string
}

func GetHttpTransport() *HttpTransport {
	return &HttpTransport{&http.Client{Timeout: time.Second * 300}, "https://api.vk.com/method/"}
}

func (t *HttpTransport) call(method Method, params Params) ([]byte, error, TransportExternalData) {
	path := t.endpoint + string(method)
	resp, err := t.client.PostForm(path, params.toMap())
	if err != nil {
		return nil, err, TransportExternalData{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err, TransportExternalData(resp.Header)
}