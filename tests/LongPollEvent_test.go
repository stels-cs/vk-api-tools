package VkApiTest

import (
	"encoding/json"
	"github.com/stels-cs/vk-api-tools"
	"reflect"
	"strconv"
	"testing"
)

func TestFillNewMessageEvent(t *testing.T) {

	data := []byte(`{"ts":1839788754,"pts":10000013,"updates":[[4,7,1,19039187,1512498938,"Привет",{"title":" ... "}]]}`)

	lpChunk := VkApi.LongPollResponse{}
	err := json.Unmarshal(data, &lpChunk)
	if err != nil {
		t.Error(err.Error())
	}
	updates := lpChunk.Updates
	if len(updates) != 1 {
		t.Error("Updates has invalid length, must be 1 but " + strconv.Itoa(len(updates)))
	}
	update := updates[0]
	if typeId, ok := update[0].(float64); ok {
		if typeId != 4 {
			t.Errorf("Invalid event type, must be 4 but %d", typeId)
		}
		msg := VkApi.MessageEvent{}
		err := msg.Fill(update)
		if err != nil {
			t.Error(err.Error())
		}
		if msg.Text != "Привет" {
			t.Error("Inavlid message text, must be Привет but " + msg.Text)
		}
		if msg.PeerId != 19039187 {
			t.Error("Inavlid message peer id, must be 19039187 but " + strconv.Itoa(msg.PeerId))
		}
	} else {
		typeName := reflect.TypeOf(update[0])
		t.Errorf("Cant parse event type, its not int! type %s", typeName)
	}
}

func TestFillMessageWithInviteEvent(t *testing.T) {

	data := []byte(`{"ts":1757798070,"pts":10000080,"updates":[[4,41,532481,2000000001,1512524680,"",{"source_act":"chat_invite_user","source_mid":"460552514","from":"19039187"}]]}`)

	lpChunk := VkApi.LongPollResponse{}
	err := json.Unmarshal(data, &lpChunk)
	if err != nil {
		t.Error(err.Error())
	}
	updates := lpChunk.Updates
	if len(updates) != 1 {
		t.Error("Updates has invalid length, must be 1 but " + strconv.Itoa(len(updates)))
	}
	update := updates[0]
	if typeId, ok := update[0].(float64); ok {
		if typeId != 4 {
			t.Errorf("Invalid event type, must be 4 but %d", typeId)
		}
		msg := VkApi.MessageEvent{}
		err := msg.Fill(update)
		if err != nil {
			t.Error(err.Error())
		}
		if msg.Text != "" {
			t.Error("Inavlid message text, must be <empty> but " + msg.Text)
		}
		if msg.PeerId != 2000000001 {
			t.Error("Inavlid message peer id, must be 2000000001 but " + strconv.Itoa(msg.PeerId))
		}
		if msg.ChatInviteUser == nil {
			t.Error("Inavlid message invite mus be from 19039187 to 460552514 but is null")
		} else if msg.ChatInviteUser.Actor != 19039187 {
			t.Errorf("Inavlid message invite actor must be 19039187 but %d", msg.ChatInviteUser.Actor)

		} else if msg.ChatInviteUser.User != 460552514 {
			t.Errorf("Inavlid message invite user must be 460552514 but %d", msg.ChatInviteUser.User)
		}
	} else {
		typeName := reflect.TypeOf(update[0])
		t.Errorf("Cant parse event type, its not int! type %s", typeName)
	}
}
