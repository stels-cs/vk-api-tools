package VkApi

import "encoding/json"

type CallbackEvent struct {
	Type    string          `json:"type"`
	Object  json.RawMessage `json:"object"`
	GroupId int             `json:"group_id"`
}

func (ev *CallbackEvent) IsMessage() bool {
	if ev.Type == "message_new" {
		return true
	}
	if ev.Type == "message_reply" {
		return true
	}
	if ev.Type == "message_edit" {
		return true
	}
	return false
}

func (ev *CallbackEvent) GetMessage() (*CallbackMessage, error) {
	m := CallbackMessage{}
	err := json.Unmarshal(ev.Object, &m)
	return &m, err
}

type CallbackEventListener func(event *CallbackEvent)

type CallbackMessage struct {
	Id        int    `json:"id"`         //
	UserId    int    `json:"user_id"`    //
	FromId    int    `json:"from_id"`    //
	Date      int    `json:"date"`       //
	ReadState int    `json:"read_state"` //, [0,1]
	Out       int    `json:"out"`        //, [0,1]
	Title     string `json:"title"`      //
	Body      string `json:"body"`       //
	Geo       struct {
		Type        string `json:"type"`
		Coordinates string `json:"coordinates"`
		Place       struct {
			Id        int     `json:"id"`
			Title     string  `json:"title"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Created   int     `json:"created"`
			Icon      string  `json:"icon"`
			Country   string  `json:"country"`
			City      string  `json:"city"`
		} `json:"place"`
	} `json:"geo"` //
	Attachments []struct {
		Type string `json:"type"`
	} `json:"attachments"` //
	FwdMessages []MessageEvent `json:"fwd_messages"` //
	Emoji       int            `json:"emoji"`        //, [0,1]
	Important   bool           `json:"important"`    //, [0,1]
	Deleted     int            `json:"deleted"`      //, [0,1]
	RandomId    int            `json:"random_id"`    //
	Payload     string         `json:"payload"`
}
