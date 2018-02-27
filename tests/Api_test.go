package VkApiTest

import (
	"github.com/stels-cs/vk-api-tools"
	"testing"
)

func ft(s string) *VkApi.FakeTransport {
	return &VkApi.FakeTransport{
		[]byte(s),
		nil,
		VkApi.TransportExternalData{},
		0,
	}
}

func TestApiUsersGet(t *testing.T) {
	tr := ft(`{"response":[{"id":1}]}`)
	api := VkApi.CreateApi("", "5.71", tr, 30)
	data, err := api.Call("users.get", VkApi.P{"user_ids": "1"})
	if err != nil {
		t.Error(err)
	}
	u, err := data.FirstAny()
	if err != nil {
		t.Error(err)
	}
	if id, _ := u.GetInt("id"); id != 1 {
		t.Errorf("Incorrect response user id must be 1 but got %d", id)
	}
}

func TestApiUsersGetWithError(t *testing.T) {
	tr := ft(`{"error":{"error_code": 113,"error_msg": "Invalid user id","request_params": []}}`)
	api := VkApi.CreateApi("", "5.71", tr, 30)
	user, err := api.Call("users.get", VkApi.P{"user_ids": "-10"})
	if err == nil {
		t.Error("Expected error, got nil")
	}
	_, err = user.Slice()
	if err == nil {
		t.Errorf("Incorrect response length for users.get user_ids=-1 must be 0 but got slice ")
	}
}

func TestApiMessagesSendWithError(t *testing.T) {
	tr := ft(`{"error":{"error_code":5}}`)
	api := VkApi.CreateApi("", "5.71", tr, 30)

	_, err := api.Call("messages.send", VkApi.P{"peer_id": "0", "messages": "Test"})
	if err == nil {
		t.Error("Messages call success but expect error")
	}
	if apiError, ok := err.(*VkApi.ApiError); ok {
		if apiError.Code != VkApi.AuthError {
			t.Errorf("Bad error code expext 5 (Vk.ApiAuthError) but %d found", apiError.Code)
		}
	} else {
		t.Errorf("Got error but is not ApiError: %s", err.Error())
	}
}
