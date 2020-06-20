package VkApi

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

type PhotoAttach struct {
	Id      int    `json:"id"`
	AlbumId int    `json:"album_id"`
	OwnerId int    `json:"owner_id"`
	UserId  int    `json:"user_id"`
	Text    string `json:"text"`
	Date    int    `json:"date"`
	Sizes   []struct {
		Type   string      `json:"type"`
		Url    string      `json:"url"`
		Width  json.Number `json:"width"`
		Height json.Number `json:"height"`
	} `json:"sizes"`
}

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
	GroupId   int
	Id        int    `json:"id"`         //
	UserId    int    `json:"user_id"`    //
	FromId    int    `json:"from_id"`    //
	PeerId    int    `json:"peer_id"`    //
	Date      int    `json:"date"`       //
	ReadState int    `json:"read_state"` //, [0,1]
	Out       int    `json:"out"`        //, [0,1]
	Title     string `json:"title"`      //
	Body      string `json:"body"`       // OLD
	Text      string `json:"text"`       //
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
		Type  string      `json:"type"`
		Photo PhotoAttach `json:"photo"`
	} `json:"attachments"` //
	FwdMessages []MessageEvent `json:"fwd_messages"` //
	Emoji       int            `json:"emoji"`        //, [0,1]
	Important   bool           `json:"important"`    //, [0,1]
	Deleted     int            `json:"deleted"`      //, [0,1]
	RandomId    int            `json:"random_id"`    //
	Payload     string         `json:"payload"`
}

func (message *CallbackMessage) HasMention() bool {
	matched, err := regexp.MatchString("\\[club"+strconv.Itoa(message.GroupId)+"\\|.*?\\]", message.Text)
	if err != nil {
		return false
	}
	return matched
}

func (message *CallbackMessage) GetAttachTypes() []string {
	ttt := make([]string, len(message.Attachments))
	for _, t := range message.Attachments {
		ttt = append(ttt, t.Type)
	}
	return ttt
}

func (message *CallbackMessage) GetTextWithoutMention() string {
	r := regexp.MustCompile("\\[club" + strconv.Itoa(message.GroupId) + "\\|.*?\\]")
	return strings.TrimSpace(r.ReplaceAllString(message.Text, ""))
}

func (message *CallbackMessage) FromUser() bool {
	return message.FromId > 0
}

func (message *CallbackMessage) FromGroup() bool {
	return message.FromId < 0
}

func (message *CallbackMessage) FromChat() bool {
	return message.PeerId >= 2e9
}

func (message *CallbackMessage) HasPayload() bool {
	return message.Payload != ""
}
