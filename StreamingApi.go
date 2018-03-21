package VkApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type StreamingCallback func(code int, event *StreamingEvent, message *StreamingServiceMessage)

type StreamingClient struct {
	Endpoint string `json:"endpoint"`
	Key      string `json:"key"`
	Http     *http.Client
	done     chan struct{}
	conn     *websocket.Conn
}

type StreamingResponse struct {
	Code  int            `json:"code"`
	Error StreamingError `json:"error"`
}

// See: https://vk.com/dev/streaming_api_docs_2?f=7.%2B%D0%A7%D1%82%D0%B5%D0%BD%D0%B8%D0%B5%2B%D0%BF%D0%BE%D1%82%D0%BE%D0%BA%D0%B0
type StreamingEvent struct {
	EventType string `json:"event_type"`
	EventId   struct {
		PostOwnerId  int `json:"post_owner_id"`
		PostId       int `json:"post_id"`
		CommentId    int `json:"comment_id"`
		SharedPostId int `json:"shared_post_id"`
	} `json:"event_id"`
	EventUrl               string   `json:"event_url"`
	Text                   string   `json:"text"`
	Action                 string   `json:"action"`
	ActionTime             int      `json:"action_time"`
	CreationTime           int      `json:"creation_time"`
	SharedPostText         string   `json:"shared_post_text"`
	SharedPostCreationTime int      `json:"shared_post_creation_time"`
	SignerId               int      `json:"signer_id"`
	Tags                   []string `json:"tags"`
	Author                 struct {
		Id                  int    `json:"id"`
		AuthorUrl           string `json:"author_url"`
		SharedPostAuthorId  int    `json:"shared_post_author_id"`
		SharedPostAuthorUrl string `json:"shared_post_author_url"`
		Platform            int    `json:"platform"`
	} `json:"author"`
}

func (se *StreamingEvent) IsShare() bool {
	return se.EventType == "share"
}

func (se *StreamingEvent) IsPost() bool {
	return se.EventType == "post"
}

func (se *StreamingEvent) IsComment() bool {
	return se.EventType == "comment"
}

func (se *StreamingEvent) IsNew() bool {
	return se.Action == "new"
}

func (se *StreamingEvent) IsUpdate() bool {
	return se.Action == "update"
}

func (se *StreamingEvent) IsDelete() bool {
	return se.Action == "delete"
}

func (se *StreamingEvent) IsRestore() bool {
	return se.Action == "restore"
}

type StreamingServiceMessage struct {
	Message     string `json:"message"`
	ServiceCode int    `json:"service_code"`
}

type StreamingError struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
}

func (e *StreamingError) Error() string {
	if e.Message == "" {
		return "VkApi.StreamingError: Unknown error"
	}
	return "VkApi.StreamingError: #" + strconv.Itoa(e.ErrorCode) + ": " + e.Message
}

type StreamingRule struct {
	Value string `json:"value"`
	Tag   string `json:"tag"`
}

func CreateStreamingClient(accessToken string) (*StreamingClient, error) {
	sc := &StreamingClient{}
	err := Exec("streaming.getServerUrl", P{"access_token": accessToken, "v": "5.73"}, sc)
	if err != nil {
		return nil, err
	} else {
		return sc, nil
	}
}

func (c *StreamingClient) call(path string, data interface{}, httpMethod string, dist interface{}) error {
	if c.Http == nil {
		c.Http = &http.Client{Timeout: time.Second * 300}
	}

	b, e := json.Marshal(data)
	if e != nil {
		return e
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.Endpoint,
		Path:     path,
		RawQuery: "key=" + c.Key,
	}

	r, err := http.NewRequest(httpMethod, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := c.Http.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, dist)
	if err != nil {
		return err
	}

	return nil
}

func (c *StreamingClient) GetRules() ([]StreamingRule, error) {
	res := struct {
		StreamingResponse
		Rules []StreamingRule `json:"rules"`
	}{}
	err := c.call("/rules", "", http.MethodGet, &res)
	if err != nil {
		return res.Rules, err
	}
	if res.Code != http.StatusOK {
		return res.Rules, &res.Error
	}
	return res.Rules, nil
}

func (c *StreamingClient) AddRule(rule StreamingRule) error {
	payload := struct {
		Rule StreamingRule `json:"rule"`
	}{
		Rule: rule,
	}
	res := &StreamingResponse{}
	err := c.call("/rules", payload, http.MethodPost, res)
	if err != nil {
		return err
	}
	if res.Code != http.StatusOK {
		return &res.Error
	}
	return nil
}

func (c *StreamingClient) DeleteRule(tag string) error {
	payload := struct {
		Tag string `json:"tag"`
	}{
		Tag: tag,
	}
	res := &StreamingResponse{}
	err := c.call("/rules", payload, "DELETE", res)
	if err != nil {
		return err
	}
	if res.Code != http.StatusOK {
		return &res.Error
	}
	return nil
}

func (c *StreamingClient) ClearAllRules() error {
	rules, err := c.GetRules()
	if err != nil {
		return err
	}
	for _, rule := range rules {
		err := c.DeleteRule(rule.Tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *StreamingClient) Start(callback StreamingCallback) error {
	u := url.URL{
		Host:     c.Endpoint,
		Scheme:   "wss",
		Path:     "/stream",
		RawQuery: "key=" + c.Key,
	}

	conn, response, err := websocket.DefaultDialer.Dial(u.String(), nil)
	c.conn = conn
	if err != nil {
		body, e := ioutil.ReadAll(response.Body)
		if e == nil {
			response.Body.Close()
			return errors.New(string(body))
		}
		return err
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			return nil
		}
		if err != nil {
			return err
		}
		buff := struct {
			Code           int                     `json:"code"`
			Event          StreamingEvent          `json:"event"`
			ServiceMessage StreamingServiceMessage `json:"service_message"`
		}{}
		err = json.Unmarshal(message, &buff)
		if err != nil {
			return err
		}
		if buff.Code == 100 {
			callback(buff.Code, &buff.Event, nil)
		} else if buff.Code == 100 {
			callback(buff.Code, nil, &buff.ServiceMessage)
		} else {
			return errors.New("VkApi.StreamingError: Unknown event code " + strconv.Itoa(buff.Code))
		}
	}

}

func (c *StreamingClient) Stop() {
	if c.conn == nil {
		return
	}
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
