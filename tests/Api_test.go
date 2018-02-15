package Vk

import (
	"testing"
	"github.com/stels-cs/vk-api-tools"
)

func TestApiUsersGet(t *testing.T) {
	api := Vk.GetApi(Vk.AccessToken{}, Vk.GetHttpTransport(), nil)
	user, err := api.Users.GetByIds([]int{1})
	if err != nil {
		t.Error(err)
	}
	if len(user) != 1 {
		t.Errorf("Incorrect response length for users.get user_ids=1 must be 1 but %d got", len(user))
	}
	if user[0].Id != 1 {
		t.Errorf("Incorrect response user id must be 1 but got %d", user[0].Id)
	}
}

func TestApiUsersGetWithError(t *testing.T) {
	api := Vk.GetApi(Vk.AccessToken{}, Vk.GetHttpTransport(), nil)
	user, err := api.Users.GetByIds([]int{-10})
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if len(user) != 0 {
		t.Errorf("Incorrect response length for users.get user_ids=-1 must be 0 but %d got", len(user))
	}
}


func TestApiMessagesSendWithError(t *testing.T) {
	api := Vk.GetApi(Vk.AccessToken{}, Vk.GetHttpTransport(), nil)

	_, err := api.Call( Vk.GetApiMethod("messages.send", Vk.Params{"peer_id":"0", "messages":"Test"}) )
	if err == nil {
		t.Error("Messages call success but expect error")
	}
	if apiError, ok := err.(*Vk.ApiError); ok {
		if apiError.Code != Vk.ApiAuthError {
			t.Errorf("Bad error code expext 5 (Vk.ApiAuthError) but %d found", apiError.Code)
		}
	} else {
		t.Errorf("Got error but is not ApiError: %s", err.Error())
	}
}