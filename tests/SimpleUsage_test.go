package VkApiTest

import (
	"github.com/stels-cs/vk-api-tools"
	"os"
	"testing"
)

var envTokenTag = "VK_ACCESS_TOKEN"

func getToken() string {
	token := os.Getenv(envTokenTag)
	return token
}

func skip(t *testing.T) {
	t.Skip("No access_token passed, pass token in enviroment variable", envTokenTag, "to run this test")
}

func TestSimpleUsage1(t *testing.T) {
	if token := getToken(); token != "" {
		res, err := VkApi.Call("users.get", VkApi.P{"user_ids": "stelscs", "access_token": token})
		if err != nil {
			t.Error(err)
		}
		user, _ := res.FirstAny()
		if i, _ := user.GetInt("id"); i != 19039187 {
			t.Error("Expected 19039187, got", i)
		}
		if i, _ := user.GetString("first_name"); i != "Иван" {
			t.Error("Expected Иван, got", i)
		}
	} else {
		skip(t)
	}
}

func TestSimpleUsage2(t *testing.T) {
	if token := getToken(); token != "" {
		res, err := VkApi.Call("users.get", VkApi.P{"user_ids": "stelscs", "access_token": token})
		if err != nil {
			t.Error(err)
		}
		user, _ := res.FirstAny()
		if i, _ := user.GetInt("id"); i != 19039187 {
			t.Error("Expected 19039187, got", i)
		}
		if i, _ := user.GetString("first_name"); i != "Иван" {
			t.Error("Expected Иван, got", i)
		}
	} else {
		skip(t)
	}
}

func TestCheckingError(t *testing.T) {
	if token := getToken(); token != "" {
		_, err := VkApi.Call("messages.send", VkApi.P{"user_ids": "stelscs", "access_token": token})
		if err == nil {
			t.Error("Expect error, but not")
		}

		if VkApi.IsApiError(err) == false {
			t.Error("Expect ApiError, but this is anoter error")
		}

		if VkApi.IsTransportError(err) == true {
			t.Error("Expect ApiError, got transport error")
		}

		e := VkApi.CastToApiError(err)

		if e.Code != VkApi.BadApiKeyError {
			t.Error("Expect Code VkApi.BadApiKeyError, got", e.Code)
		}
	} else {
		skip(t)
	}
}

func TestCheckTransportError(t *testing.T) {
	tr := ft(`LOL KEK CHEBYREK`)
	api := VkApi.CreateApi("", "", tr, 0)

	_, err := api.Call("users.get", VkApi.P{"user_ids": "1"})

	if VkApi.IsTransportError(err) == false {
		t.Error("Expect TransportError, got anoter error")
	}

	if VkApi.IsApiError(err) == true {
		t.Error("Expect TransportError, got ApiError")
	}

	e := VkApi.CastToTransportError(err)

	if e.Method != "users.get" {
		t.Error("Expected method users.get, got", e.Method)
	}
}
