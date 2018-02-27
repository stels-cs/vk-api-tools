package VkApi

import (
	"io/ioutil"
	"net/http"
	"time"
)

type HttpTransport struct {
	client   *http.Client
	endpoint string
}

// Возвращает транспорт по умолчанию для передачи его в VkApi.CreateApi
func GetHttpTransport() *HttpTransport {
	return &HttpTransport{&http.Client{Timeout: time.Second * 300}, "https://api.vk.com/method/"}
}

func (t *HttpTransport) call(method string, params P) ([]byte, error, TransportExternalData) {
	path := t.endpoint + string(method)
	resp, err := t.client.PostForm(path, params.toMap())
	if err != nil {
		return nil, err, TransportExternalData{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err, TransportExternalData(resp.Header)
}
