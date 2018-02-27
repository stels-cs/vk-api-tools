package VkApi

import (
	"strings"
	"time"
)

type TransportExternalData map[string][]string

func (d *TransportExternalData) toString() string {
	str := ""
	for k, v := range *d {
		if str != "" {
			str += "\n"
		}
		str += k + "=" + strings.Join(v, ",")
	}
	return str
}

type Transport interface {
	call(method string, params P) ([]byte, error, TransportExternalData)
}

type TransportError struct {
	Method      string
	Params      P
	Response    []byte
	Headers     TransportExternalData
	ParentError error
}

func (e *TransportError) Error() string {
	if e.ParentError != nil {
		return e.ParentError.Error()
	} else {
		return e.DebugInfo()
	}
}

func (e *TransportError) DebugInfo() string {
	if len(e.Response) <= 1000 {
		return "TransportError: " + callToString(e.Method, e.Params) + "\n" + e.Headers.toString() + "\n" + string(e.Response)
	} else {
		startIndex := len(e.Response) - 1000
		s := string(e.Response[:1000]) + "..." + string(e.Response[startIndex:])
		return "TransportError: " + callToString(e.Method, e.Params) + "\n" + e.Headers.toString() + "\n" + s
	}
}

type FakeTransport struct {
	Response     []byte
	Err          error
	ExternalData TransportExternalData
	SleepTime    int64
}

func (t *FakeTransport) call(method string, params P) ([]byte, error, TransportExternalData) {
	if t.SleepTime > 0 {
		time.Sleep(time.Duration(t.SleepTime) * time.Millisecond)
	}
	return t.Response, t.Err, t.ExternalData
}

type FakeTransportPoll struct {
	Data []FakeTransport
}

func (t *FakeTransportPoll) call(method string, params P) ([]byte, error, TransportExternalData) {
	tr := t.Data[0]
	t.Data = t.Data[1:]
	return tr.call(method, params)
}
