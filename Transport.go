package Vk

import "strings"

type TransportExternalData map[string][]string

func (d *TransportExternalData) toString() string {
	str := ""
	for k,v:= range *d {
		if str != "" {str += "\n"}
		str += k + "=" + strings.Join(v, ",")
	}
	return str
}

type Transport interface {
	call(method Method, params Params) ([]byte, error, TransportExternalData)
}

type TransportBadResponse struct {
	method Method
	params Params
	response []byte
	headers TransportExternalData
}

func (e *TransportBadResponse) Error() string {
	if len(e.response) <= 1000 {
		return "BadResponse\n" + callToString(e.method, e.params) + "\n\n" + e.headers.toString() + "\n" + string(e.response)
	} else {
		return "BadResponse Large response \n" + callToString(e.method, e.params) + "\n\n" + e.headers.toString() + "\n" + string(e.response[:1000])
	}
}